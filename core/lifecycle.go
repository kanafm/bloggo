package core

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

type buildPlan struct {
	jobs []buildJob
}

type buildJob struct {
	entryPoint string
	inputFile  string
	outputFile string
	processor  Processor
}

// Build plans and executes a site build.
func Build(request BuildRequest) error {
	plan, err := createBuildPlan(request)
	if err != nil {
		return err
	}

	return executeBuildPlan(plan)
}

func createBuildPlan(request BuildRequest) (buildPlan, error) {
	plan := buildPlan{}
	router := request.Router
	if router == nil {
		router = MirrorRouter{}
	}

	absoluteOutDir, err := filepath.Abs(request.OutDir)
	if err != nil {
		return plan, fmt.Errorf("resolve output directory %q: %w", request.OutDir, err)
	}

	claimedInputs := make(map[string]string)
	claimedOutputs := make(map[string]string)

	for _, entryPoint := range request.EntryPoints {
		absoluteEntryPoint, err := filepath.Abs(entryPoint)
		if err != nil {
			return plan, fmt.Errorf("resolve entry point %q: %w", entryPoint, err)
		}

		if pathContains(absoluteEntryPoint, absoluteOutDir) {
			return plan, fmt.Errorf("output directory %q is inside entry point %q", absoluteOutDir, absoluteEntryPoint)
		}

		err = filepath.WalkDir(entryPoint, func(path string, entry fs.DirEntry, walkErr error) error {
			if walkErr != nil {
				return walkErr
			}
			if entry.IsDir() {
				return nil
			}

			absoluteInput, err := filepath.Abs(path)
			if err != nil {
				return fmt.Errorf("resolve input file %q: %w", path, err)
			}
			if previousEntryPoint, exists := claimedInputs[absoluteInput]; exists {
				return fmt.Errorf(
					"input file %q is included by entry points %q and %q",
					absoluteInput,
					previousEntryPoint,
					entryPoint,
				)
			}
			claimedInputs[absoluteInput] = entryPoint

			relativeInput, err := filepath.Rel(entryPoint, path)
			if err != nil {
				return fmt.Errorf("route input file %q: %w", path, err)
			}

			relativeOutput, ok := router.Route(relativeInput)
			if !ok {
				fmt.Println(". skipped: router did not generate a route for input", path)
				return nil
			}

			outputFile, err := resolveOutputFile(absoluteOutDir, path, relativeOutput)
			if err != nil {
				return err
			}

			processor, err := selectProcessor(path, request.Processors)
			if err != nil {
				return err
			}

			if previousInput, exists := claimedOutputs[outputFile]; exists {
				return fmt.Errorf(
					"output collision for %q: inputs %q and %q",
					outputFile,
					previousInput,
					path,
				)
			}

			claimedOutputs[outputFile] = path
			plan.jobs = append(plan.jobs, buildJob{
				entryPoint: entryPoint,
				inputFile:  path,
				outputFile: outputFile,
				processor:  processor,
			})

			return nil
		})
		if err != nil {
			return plan, fmt.Errorf("plan entry point %q: %w", entryPoint, err)
		}
	}

	sort.Slice(plan.jobs, func(i, j int) bool {
		if plan.jobs[i].outputFile == plan.jobs[j].outputFile {
			return plan.jobs[i].inputFile < plan.jobs[j].inputFile
		}
		return plan.jobs[i].outputFile < plan.jobs[j].outputFile
	})

	return plan, nil
}

func resolveOutputFile(outDir string, inputFile string, route string) (string, error) {
	if route == "" {
		return "", fmt.Errorf("router generated an empty route for input %q", inputFile)
	}
	if filepath.IsAbs(route) || filepath.VolumeName(route) != "" {
		return "", fmt.Errorf("router generated an absolute route %q for input %q", route, inputFile)
	}

	cleanRoute := filepath.Clean(route)
	if cleanRoute == "." || cleanRoute == ".." || strings.HasPrefix(cleanRoute, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("router generated an invalid route %q for input %q", route, inputFile)
	}

	outputFile := filepath.Join(outDir, cleanRoute)
	if !pathContains(outDir, outputFile) {
		return "", fmt.Errorf("router generated a route outside the output directory %q for input %q", route, inputFile)
	}

	return outputFile, nil
}

func executeBuildPlan(plan buildPlan) error {
	for _, job := range plan.jobs {
		fmt.Println(". processing ", job.inputFile)

		if err := os.MkdirAll(filepath.Dir(job.outputFile), 0o755); err != nil {
			return fmt.Errorf("create output directory for %q: %w", job.outputFile, err)
		}
		if err := job.processor.Handle(job.inputFile, job.outputFile); err != nil {
			return fmt.Errorf("process %q into %q: %w", job.inputFile, job.outputFile, err)
		}

		fmt.Println("\t. -> ", job.outputFile)
	}

	return nil
}

func selectProcessor(inputFile string, processors []Processor) (Processor, error) {
	for _, processor := range processors {
		if processor.CanHandle(inputFile) {
			return processor, nil
		}
	}

	return nil, fmt.Errorf("no processor can handle %q", inputFile)
}

func pathContains(parent string, child string) bool {
	relative, err := filepath.Rel(parent, child)
	if err != nil {
		return false
	}

	return relative == "." || (relative != ".." && !strings.HasPrefix(relative, ".."+string(filepath.Separator)))
}
