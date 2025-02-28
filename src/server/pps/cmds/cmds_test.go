// TODO(msteffen) Add tests for:
//
// - restart datum
// - stop job
// - delete job
//
// - inspect job
// - list job
//
// - create pipeline
// - create pipeline --push-images (re-enable existing test)
// - update pipeline
// - delete pipeline
//
// - inspect pipeline
// - list pipeline
//
// - start pipeline
// - stop pipeline
//
// - list datum
// - inspect datum
// - logs

package cmds

import (
	"bytes"
	"fmt"
	"os"
	"testing"

	"github.com/pachyderm/pachyderm/v2/src/internal/require"
	tu "github.com/pachyderm/pachyderm/v2/src/internal/testutil"
)

const badJSON1 = `
{
"356weryt

}
`

const badJSON2 = `{
{
    "a": 1,
    "b": [23,4,4,64,56,36,7456,7],
    "c": {"e,f,g,h,j,j},
    "d": 3452.36456,
}
`

func TestSyntaxErrorsReportedCreatePipeline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	require.NoError(t, tu.BashCmd(`
		echo -n '{{.badJSON1}}' \
		  | ( pachctl create pipeline -f - 2>&1 || true ) \
		  | match "malformed pipeline spec"

		echo -n '{{.badJSON2}}' \
		  | ( pachctl create pipeline -f - 2>&1 || true ) \
		  | match "malformed pipeline spec"
		`,
		"badJSON1", badJSON1,
		"badJSON2", badJSON2,
	).Run())
}

func TestSyntaxErrorsReportedUpdatePipeline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	require.NoError(t, tu.BashCmd(`
		echo -n '{{.badJSON1}}' \
		  | ( pachctl update pipeline -f - 2>&1 || true ) \
		  | match "malformed pipeline spec"

		echo -n '{{.badJSON2}}' \
		  | ( pachctl update pipeline -f - 2>&1 || true ) \
		  | match "malformed pipeline spec"
		`,
		"badJSON1", badJSON1,
		"badJSON2", badJSON2,
	).Run())
}

func TestRawFullPipelineInfo(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	require.NoError(t, tu.BashCmd(`
		yes | pachctl delete all
	`).Run())
	require.NoError(t, tu.BashCmd(`
		pachctl create repo data
		pachctl put file data@master:/file <<<"This is a test"
		pachctl create pipeline <<EOF
		  {
		    "pipeline": {"name": "{{.pipeline}}"},
		    "input": {
		      "pfs": {
		        "glob": "/*",
		        "repo": "data"
		      }
		    },
		    "transform": {
		      "cmd": ["bash"],
		      "stdin": ["cp /pfs/data/file /pfs/out"]
		    }
		  }
		EOF
		`,
		"pipeline", tu.UniqueString("p-")).Run())
	require.NoError(t, tu.BashCmd(`
		pachctl flush commit data@master

		# make sure the results have the full pipeline info, including version
		pachctl list job --raw --history=all \
			| match "pipeline_version"
		`).Run())
}

// TestJSONMultiplePipelines tests that pipeline specs with multiple pipelines
// in them continue to be accepted by 'pachctl create pipeline'. We may want to
// stop supporting this behavior eventually, but Pachyderm has supported it
// historically, so we should continue to support it until we formally deprecate
// it.
func TestJSONMultiplePipelines(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	require.NoError(t, tu.BashCmd(`
		yes | pachctl delete all
		pachctl create repo input
		pachctl create pipeline -f - <<EOF
		{
		  "pipeline": {
		    "name": "first"
		  },
		  "input": {
		    "pfs": {
		      "glob": "/*",
		      "repo": "input"
		    }
		  },
		  "transform": {
		    "cmd": [ "/bin/bash" ],
		    "stdin": [
		      "cp /pfs/input/* /pfs/out"
		    ]
		  }
		}
		{
		  "pipeline": {
		    "name": "second"
		  },
		  "input": {
		    "pfs": {
		      "glob": "/*",
		      "repo": "first"
		    }
		  },
		  "transform": {
		    "cmd": [ "/bin/bash" ],
		    "stdin": [
		      "cp /pfs/first/* /pfs/out"
		    ]
		  }
		}
		EOF

		pachctl start commit input@master
		echo foo | pachctl put file input@master:/foo
		echo bar | pachctl put file input@master:/bar
		echo baz | pachctl put file input@master:/baz
		pachctl finish commit input@master
		pachctl flush commit input@master
		pachctl get file second@master:/foo | match foo
		pachctl get file second@master:/bar | match bar
		pachctl get file second@master:/baz | match baz
		`,
	).Run())
	require.NoError(t, tu.BashCmd(`pachctl list pipeline`).Run())
}

// TestJSONStringifiedNumberstests that JSON pipelines may use strings to
// specify numeric values such as a pipeline's parallelism (a feature of gogo's
// JSON parser).
func TestJSONStringifiedNumbers(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	require.NoError(t, tu.BashCmd(`
		yes | pachctl delete all
		pachctl create repo input
		pachctl create pipeline -f - <<EOF
		{
		  "pipeline": {
		    "name": "first"
		  },
		  "input": {
		    "pfs": {
		      "glob": "/*",
		      "repo": "input"
		    }
		  },
		  "parallelism_spec": {
		    "constant": "1"
		  },
		  "transform": {
		    "cmd": [ "/bin/bash" ],
		    "stdin": [
		      "cp /pfs/input/* /pfs/out"
		    ]
		  }
		}
		EOF

		pachctl start commit input@master
		echo foo | pachctl put file input@master:/foo
		echo bar | pachctl put file input@master:/bar
		echo baz | pachctl put file input@master:/baz
		pachctl finish commit input@master
		pachctl flush commit input@master
		pachctl get file first@master:/foo | match foo
		pachctl get file first@master:/bar | match bar
		pachctl get file first@master:/baz | match baz
		`,
	).Run())
	require.NoError(t, tu.BashCmd(`pachctl list pipeline`).Run())
}

// TestJSONMultiplePipelinesError tests that when creating multiple pipelines
// (which only the encoding/json parser can parse) you get an error indicating
// the problem in the JSON, rather than an error complaining about multiple
// documents.
func TestJSONMultiplePipelinesError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	// pipeline spec has no quotes around "name" in first pipeline
	require.NoError(t, tu.BashCmd(`
		yes | pachctl delete all
		pachctl create repo input
		( pachctl create pipeline -f - 2>&1 <<EOF || true
		{
		  "pipeline": {
		    name: "first"
		  },
		  "input": {
		    "pfs": {
		      "glob": "/*",
		      "repo": "input"
		    }
		  },
		  "transform": {
		    "cmd": [ "/bin/bash" ],
		    "stdin": [
		      "cp /pfs/input/* /pfs/out"
		    ]
		  }
		}
		{
		  "pipeline": {
		    "name": "second"
		  },
		  "input": {
		    "pfs": {
		      "glob": "/*",
		      "repo": "first"
		    }
		  },
		  "transform": {
		    "cmd": [ "/bin/bash" ],
		    "stdin": [
		      "cp /pfs/first/* /pfs/out"
		    ]
		  }
		}
		EOF
		) | match "invalid character 'n' looking for beginning of object key string"
		`,
	).Run())
}

func TestRunPipeline(t *testing.T) {
	// TODO(2.0 optional): Run pipeline creates a dangling commit, but uses the stats branch for stats commits.
	// Since there is no relationship between the commits created by run pipeline, there should be no
	// relationship between the stats commits. So, the stats commits should also be dangling commits.
	// This might be easier to address when global IDs are implemented.
	t.Skip("Run pipeline does not work correctly with stats enabled")
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	pipeline := tu.UniqueString("p-")
	require.NoError(t, tu.BashCmd(`
		yes | pachctl delete all
	`).Run())
	require.NoError(t, tu.BashCmd(`
		pachctl create repo data

		# Create a commit and put some data in it
		commit1=$(pachctl start commit data@master)
		echo "file contents" | pachctl put file data@${commit1}:/file -f -
		pachctl finish commit data@${commit1}

		pachctl put file data@master:/file <<<"This is a test"
		pachctl create pipeline <<EOF
			{
			  "pipeline": {"name": "{{.pipeline}}"},
			  "input": {
			    "pfs": {
			      "glob": "/*",
			      "repo": "data"
			    }
			  },
			  "transform": {
			    "cmd": ["bash"],
			    "stdin": ["cp /pfs/data/file /pfs/out"]
			  }
			}
		EOF

		# make sure the run pipeline command accepts various provenance formats
		pachctl run pipeline {{.pipeline}} data@${commit1}
		pachctl run pipeline {{.pipeline}} data@master=${commit1}
		pachctl run pipeline {{.pipeline}} data@master
		`,
		"pipeline", pipeline).Run())
}

// TestYAMLPipelineSpec tests creating a pipeline with a YAML pipeline spec
func TestYAMLPipelineSpec(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	// Note that BashCmd dedents all lines below including the YAML (which
	// wouldn't parse otherwise)
	require.NoError(t, tu.BashCmd(`
		yes | pachctl delete all
		pachctl create repo input
		pachctl create pipeline -f - <<EOF
		pipeline:
		  name: first
		input:
		  pfs:
		    glob: /*
		    repo: input
		transform:
		  cmd: [ /bin/bash ]
		  stdin:
		    - "cp /pfs/input/* /pfs/out"
		---
		pipeline:
		  name: second
		input:
		  pfs:
		    glob: /*
		    repo: first
		transform:
		  cmd: [ /bin/bash ]
		  stdin:
		    - "cp /pfs/first/* /pfs/out"
		EOF
		pachctl start commit input@master
		echo foo | pachctl put file input@master:/foo
		echo bar | pachctl put file input@master:/bar
		echo baz | pachctl put file input@master:/baz
		pachctl finish commit input@master
		pachctl flush commit input@master
		pachctl get file second@master:/foo | match foo
		pachctl get file second@master:/bar | match bar
		pachctl get file second@master:/baz | match baz
		`,
	).Run())
}

func TestListPipelineFilter(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	require.NoError(t, tu.BashCmd(`
		yes | pachctl delete all
	`).Run())
	pipeline := tu.UniqueString("pipeline")
	require.NoError(t, tu.BashCmd(`
		yes | pachctl delete all
		pachctl create repo input
		pachctl create pipeline -f - <<EOF
		{
		"pipeline": {
		  "name": "{{.pipeline}}"
		},
		"input": {
		  "pfs": {
		    "glob": "/*",
		    "repo": "input"
		  }
		},
		  "parallelism_spec": {
		    "constant": "1"
		  },
		"transform": {
		  "cmd": [ "/bin/bash" ],
		  "stdin": [
		    "cp /pfs/input/* /pfs/out"
		  ]
		}
		}
		EOF

		echo foo | pachctl put file input@master:/foo
		pachctl flush commit input@master
		# make sure we see the pipeline with the appropriate state filters
		pachctl list pipeline | match {{.pipeline}}
		pachctl list pipeline --state crashing --state failure | match -v {{.pipeline}}
		pachctl list pipeline --state starting --state running | match {{.pipeline}}
	`, "pipeline", pipeline,
	).Run())
}

// TestYAMLError tests that when creating pipelines using a YAML spec with an
// error, you get an error indicating the problem in the YAML, rather than an
// error complaining about multiple documents.
//
// Note that with the new parsing method added to support free-form fields like
// TFJob, this YAML is parsed, serialized and then re-parsed, so the error will
// refer to "json" (the format used for the canonicalized pipeline), but the
// issue referenced by the error (use of a string instead of an array for 'cmd')
// is the main problem below
func TestYAMLError(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	// "cmd" should be a list, instead of a string
	require.NoError(t, tu.BashCmd(`
		yes | pachctl delete all
		pachctl create repo input
		( pachctl create pipeline -f - 2>&1 <<EOF || true
		pipeline:
		  name: first
		input:
		  pfs:
		    glob: /*
		    repo: input
		transform:
		  cmd: /bin/bash # should be list, instead of string
		  stdin:
		    - "cp /pfs/input/* /pfs/out"
		EOF
		) | match "cannot unmarshal string into Go value of type \[\]json.RawMessage"
		`,
	).Run())
}

func TestTFJobBasic(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	require.NoError(t, tu.BashCmd(`
		yes | pachctl delete all
		pachctl create repo input
		( pachctl create pipeline -f - 2>&1 <<EOF || true
		pipeline:
		  name: first
		input:
		  pfs:
		    glob: /*
		    repo: input
		tf_job:
		  apiVersion: kubeflow.org/v1
		  kind: TFJob
		  metadata:
		    generateName: tfjob
		    namespace: kubeflow
		  spec:
		    tfReplicaSpecs:
		      PS:
		        replicas: 1
		        restartPolicy: OnFailure
		        template:
		          spec:
		            containers:
		            - name: tensorflow
		              image: gcr.io/your-project/your-image
		              command:
		                - python
		                - -m
		                - trainer.task
		                - --batch_size=32
		                - --training_steps=1000
		      Worker:
		        replicas: 3
		        restartPolicy: OnFailure
		        template:
		          spec:
		            containers:
		            - name: tensorflow
		              image: gcr.io/your-project/your-image
		              command:
		                - python
		                - -m
		                - trainer.task
		                - --batch_size=32
		                - --training_steps=1000
		EOF
		) | match "not supported yet"
		`,
	).Run())
}

// TestYAMLSecret tests creating a YAML pipeline with a secret (i.e. the fix for
// https://github.com/pachyderm/pachyderm/v2/issues/4119)
func TestYAMLSecret(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	// Note that BashCmd dedents all lines below including the YAML (which
	// wouldn't parse otherwise)
	require.NoError(t, tu.BashCmd(`
		yes | pachctl delete all

		# kubectl get secrets >&2
		kubectl delete secrets/test-yaml-secret || true
		kubectl create secret generic test-yaml-secret --from-literal=my-key=my-value

		pachctl create repo input
		pachctl put file input@master:/foo <<<"foo"
		pachctl create pipeline -f - <<EOF
		  pipeline:
		    name: pipeline
		  input:
		    pfs:
		      glob: /*
		      repo: input
		  transform:
		    cmd: [ /bin/bash ]
		    stdin:
		      - "env | grep MY_SECRET >/pfs/out/vars"
		    secrets:
		      - name: test-yaml-secret
		        env_var: MY_SECRET
		        key: my-key
		EOF
		pachctl flush commit input@master
		pachctl get file pipeline@master:/vars | match MY_SECRET=my-value
		`,
	).Run())
}

// TestYAMLTimestamp tests creating a YAML pipeline with a timestamp (i.e. the
// fix for https://github.com/pachyderm/pachyderm/v2/issues/4209)
func TestYAMLTimestamp(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	// Note that BashCmd dedents all lines below including the YAML (which
	// wouldn't parse otherwise)
	require.NoError(t, tu.BashCmd(`
		yes | pachctl delete all

		# If the pipeline comes up without error, then the YAML parsed
		pachctl create pipeline -f - <<EOF
		  pipeline:
		    name: pipeline
		  input:
		    cron:
		      name: in
		      start: "2019-10-10T22:30:05Z"
		      spec: "@yearly"
		  transform:
		    cmd: [ /bin/bash ]
		    stdin:
		      - "cp /pfs/in/* /pfs/out"
		EOF
		pachctl list pipeline | match 'pipeline'
		`,
	).Run())
}

func TestEditPipeline(t *testing.T) {
	// TODO: Edit pipeline needs to be implemented in V2.
	t.Skip("Edit pipeline not implemented in V2")
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	require.NoError(t, tu.BashCmd(`
		yes | pachctl delete all
	`).Run())
	require.NoError(t, tu.BashCmd(`
		pachctl create repo data
		pachctl create pipeline <<EOF
		  pipeline:
		    name: my-pipeline
		  input:
		    pfs:
		      glob: /*
		      repo: data
		  transform:
		    cmd: [ /bin/bash ]
		    stdin:
		      - "cp /pfs/data/* /pfs/out"
		EOF
		`).Run())
	require.NoError(t, tu.BashCmd(`
		EDITOR="cat -u" pachctl edit pipeline my-pipeline -o yaml \
		| match 'name: my-pipeline' \
		| match 'repo: data' \
		| match 'cmd:' \
		| match 'cp /pfs/data/\* /pfs/out'
		`).Run())
}

func TestPipelineBuildLifecyclePython(t *testing.T) {
	t.Skip("not implemented in V2")
	require.NoError(t, tu.BashCmd("yes | pachctl delete all").Run())
	pipeline := testPipelineBuildLifecycle(t, "python", "python")

	// the python example also contains a `.pachignore`, so we can verify it's
	// intended behavior here
	require.NoError(t, tu.BashCmd(`
		pachctl get file {{.pipeline}}_build@source:/.pachignore
	`, "pipeline", pipeline).Run())
	require.YesError(t, tu.BashCmd(`
		pachctl get file {{.pipeline}}_build@source:/foo.txt
	`, "pipeline", pipeline).Run())
}

func TestPipelineBuildLifecyclePythonNoDeps(t *testing.T) {
	t.Skip("not implemented in V2")
	require.NoError(t, tu.BashCmd("yes | pachctl delete all").Run())
	testPipelineBuildLifecycle(t, "python", "python_no_deps")
}

func TestPipelineBuildLifecycleGo(t *testing.T) {
	t.Skip("not implemented in V2")
	require.NoError(t, tu.BashCmd("yes | pachctl delete all").Run())
	testPipelineBuildLifecycle(t, "go", "go")
}

func TestAuthorizedPipelineBuildLifecycle(t *testing.T) {
	t.Skip("not implemented in V2")
	require.NoError(t, tu.BashCmd("yes | pachctl delete all").Run())
	tu.ActivateAuth(t)

	defer tu.DeleteAll(t) // make sure to clean up auth

	testPipelineBuildLifecycle(t, "go", "go")
}

//lint:ignore U1000 unreachable in v2
func testPipelineBuildLifecycle(t *testing.T, lang, dir string) string {
	t.Helper()
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	prefix := "../../../../etc/testing/pipeline-build"

	// reset and create some test input
	require.NoError(t, tu.BashCmd(`
		pachctl create repo in
		pachctl put file -r in@master:/ -f {{.prefix}}/input
	`, "prefix", prefix).Run())

	// give pipeline a unique name to work around
	// github.com/kubernetes/kubernetes/issues/82130
	pipeline := tu.UniqueString("test-pipeline-build")
	spec := fmt.Sprintf(`
		{
		  "pipeline": {
		    "name": %q
		  },
		  "transform": {
		    "build": {
		      "language": %q
		    }
		  },
		  "input": {
		    "pfs": {
		      "repo": "in",
		      "glob": "/*"
		    }
		  }
		}
	`, pipeline, lang)

	// test a barebones pipeline with a build spec and verify results
	require.NoError(t, tu.BashCmd(`
		cd {{.prefix}}/{{.dir}}
		pachctl create pipeline <<EOF
		{{.spec}}
		EOF
		pachctl flush commit in@master
		`,
		"dir", dir,
		"spec", spec,
		"prefix", prefix,
	).Run())
	require.YesError(t, tu.BashCmd(fmt.Sprintf(`
		pachctl list pipeline --state failure | match %s
	`, pipeline)).Run())
	verifyPipelineBuildOutput(t, pipeline, "0")

	// update the barebones pipeline and verify results
	require.NoError(t, tu.BashCmd(`
		cd {{.prefix}}/{{.dir}}
		pachctl update pipeline <<EOF
		{{.spec}}
		EOF
		pachctl flush commit in@master
		`,
		"dir", dir,
		"spec", spec,
		"prefix", prefix,
	).Run())
	verifyPipelineBuildOutput(t, pipeline, "0")

	// update the pipeline with a custom cmd and verify results
	spec = fmt.Sprintf(`
		{
		  "pipeline": {
		    "name": %q
		  },
		  "transform": {
		    "cmd": [
		      "sh",
		      "/pfs/build/run.sh",
		      "_"
		    ],
		    "build": {
		      "language": %q,
		      "path": "."
		    }
		  },
		  "input": {
		    "pfs": {
		      "repo": "in",
		      "glob": "/*"
		    }
		  }
		}
	`, pipeline, lang)

	require.NoError(t, tu.BashCmd(`
		cd {{.prefix}}/{{.dir}}
		pachctl update pipeline --reprocess <<EOF
		{{.spec}}
		EOF
		pachctl flush commit in@master
		`,
		"dir", dir,
		"spec", spec,
		"prefix", prefix,
	).Run())
	verifyPipelineBuildOutput(t, pipeline, "_")
	return pipeline
}

//lint:ignore U1000 unreachable in v2
func verifyPipelineBuildOutput(t *testing.T, pipeline, prefix string) {
	t.Helper()

	require.NoError(t, tu.BashCmd(`
		pachctl flush commit in@master
		pachctl get file {{.pipeline}}@master:/1.txt | match {{.prefix}}{{.prefix}}{{.prefix}}1
		pachctl get file {{.pipeline}}@master:/11.txt | match {{.prefix}}{{.prefix}}11
		pachctl get file {{.pipeline}}@master:/111.txt | match {{.prefix}}111
		`,
		"pipeline", pipeline,
		"prefix", prefix,
	).Run())
}

func TestMissingPipeline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	// should fail because there's no pipeline object in the spec
	require.YesError(t, tu.BashCmd(`
		pachctl create pipeline <<EOF
		  {
		    "transform": {
			  "image": "ubuntu:20.04"
		    },
		    "input": {
		      "pfs": {
		        "repo": "in",
		        "glob": "/*"
		      }
		    }
		  }
		EOF
	`).Run())
}

func TestUnnamedPipeline(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	// should fail because there's no pipeline name
	require.YesError(t, tu.BashCmd(`
		pachctl create pipeline <<EOF
		  {
		    "pipeline": {},
		    "transform": {
			  "image": "ubuntu:20.04"
		    },
		    "input": {
		      "pfs": {
		        "repo": "in",
		        "glob": "/*"
		      }
		    }
		  }
		EOF
	`).Run())
}

func TestPipelineBuildSpout(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	// should fail because pipeline build can't be used w/ spouts
	require.YesError(t, tu.BashCmd(`
		pachctl create pipeline <<EOF
		  {
		    "pipeline": {
		      "name": "test"
		    },
		    "transform": {
		      "build": {
		        "language": "go"
		      }
		    },
		    "spout": {}
		  }
		EOF
	`).Run())
}

func TestPipelineBuildMissingInput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	// should fail because pipeline build can't be used w/ spouts
	require.YesError(t, tu.BashCmd(`
		pachctl create pipeline <<EOF
		  {
		    "pipeline": {
		      "name": "test"
		    },
		    "transform": {
		      "build": {
		        "language": "go"
		      }
		    }
		  }
		EOF
	`).Run())
}

func TestPipelineBuildBadInput(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	// should fail because 'build' and 'source' are reserved input names when
	// using pipeline builds
	require.YesError(t, tu.BashCmd(`
		pachctl create pipeline <<EOF
		  {
		    "pipeline": {
		      "name": "test"
		    },
		    "transform": {
		      "build": {
		        "language": "go"
		      }
		    },
		    "input": {
		      "pfs": {
		        "name": "build",
		        "repo": "in",
		        "glob": "/*"
		      }
		    }
		  }
		EOF
	`).Run())
	require.YesError(t, tu.BashCmd(`
		pachctl create pipeline <<EOF
		  {
		    "pipeline": {
		      "name": "test"
		    },
		    "transform": {
		      "build": {
		        "language": "go"
		      }
		    },
		    "input": {
		      "pfs": {
		        "name": "source",
		        "repo": "in",
		        "glob": "/*"
		      }
		    }
		  }
		EOF
	`).Run())
}

func TestPipelineBuildMissingPath(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	// should fail because 'build' and 'source' are reserved input names for
	// when using pipeline builds
	require.YesError(t, tu.BashCmd(`
		pachctl create pipeline <<EOF
			{
			  "pipeline": {
			    "name": "test"
			  },
			  "transform": {
			    "build": {
			      "language": "go",
			      "path": "/some/path/that/doesnt/exist"
			    }
			  },
			  "input": {
			    "pfs": {
			      "repo": "in",
			      "glob": "/*"
			    }
			  }
			}
		EOF
	`).Run())
}

// func TestPushImages(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("Skipping integration tests in short mode")
// 	}
// 	ioutil.WriteFile("test-push-images.json", []byte(`{
//   "pipeline": {
//     "name": "test_push_images"
//   },
//   "transform": {
//     "cmd": [ "true" ],
// 	"image": "test-job-shim"
//   }
// }`), 0644)
// 	os.Args = []string{"pachctl", "create", "pipeline", "--push-images", "-f", "test-push-images.json"}
// 	require.NoError(t, rootCmd().Execute())
// }

func runPipelineWithImageGetStderr(t *testing.T, image string) (string, error) {
	// reset and create some test input
	require.NoError(t, tu.BashCmd(`
		yes | pachctl delete all
		pachctl create repo in
		pachctl put file -r in@master:/ -f ../../../../etc/testing/pipeline-build/input
	`).Run())

	cmd := tu.BashCmd(`
		pachctl create pipeline <<EOF
		  {
		    "pipeline": {
		      "name": "first"
		    },
		    "transform": {
		      "image": "{{.image}}"
		    },
		    "input": {
		      "pfs": {
		        "repo": "in",
		        "glob": "/*"
		      }
		    }
		  }
		EOF
	`, "image", image)
	buf := &bytes.Buffer{}
	cmd.Stderr = buf
	// cmd.Stdout = os.Stdout // uncomment for debugging
	cmd.Env = os.Environ()
	err := cmd.Run()
	stderr := buf.String()
	return stderr, err
}

func TestPipelineBuildRunCron(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}

	// build a cron pipeline successfully
	require.NoError(t, tu.BashCmd(`
		pachctl create pipeline <<EOF
		{
		  "pipeline": {
		    "name": "crontest"
		  },
		  "input": {
		    "cron": {
		      "name": "tick",
		      "spec": "*/1 * * * *",
		      "overwrite": true
		    }
		  },
		  "transform": {
		    "build":{
		      "language": "go"
		    },
		    "cmd":["echo", "tick"]
		  },
		  "enable_stats": true
		  }
		EOF
	`).Run())

	// and make sure you can use `run cron` on it
	require.NoError(t, tu.BashCmd(`
		pachctl run cron crontest
	`).Run())
}

func TestWarningLatestTag(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	// should emit a warning because user specified latest tag on docker image
	stderr, err := runPipelineWithImageGetStderr(t, "ubuntu:latest")
	require.NoError(t, err, "%v", err)
	require.Matches(t, "WARNING", stderr)
}

func TestWarningEmptyTag(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	// should emit a warning because user specified empty tag, equivalent to
	// :latest
	stderr, err := runPipelineWithImageGetStderr(t, "ubuntu")
	require.NoError(t, err)
	require.Matches(t, "WARNING", stderr)
}

func TestNoWarningTagSpecified(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration tests in short mode")
	}
	// should not emit a warning (stderr should be empty) because user
	// specified non-empty, non-latest tag
	stderr, err := runPipelineWithImageGetStderr(t, "ubuntu:20.04")
	require.NoError(t, err)
	require.Equal(t, "", stderr)
}
