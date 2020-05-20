package resources_test

import (
	"bytes"
	"fmt"
	"testing"
	"text/template"

	"github.com/hashicorp/terraform-plugin-sdk/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

type (
	AccTaskTestSettings struct {
		WarehouseName string

		DatabaseName string

		RootTask *TaskSettings

		ChildTask *TaskSettings

		SoloTask *TaskSettings
	}

	TaskSettings struct {
		Name     string
		Enabled  bool
		Schema   string
		SQL      string
		Schedule string
		Comment  string
		When     string
	}
)

var (
	warehousename = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	databasename  = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	rootname      = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	childname     = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)
	soloname      = acctest.RandStringFromCharSet(10, acctest.CharSetAlpha)

	initialState = &AccTaskTestSettings{
		WarehouseName: warehousename,
		DatabaseName:  databasename,

		RootTask: &TaskSettings{
			Name:     rootname,
			Schema:   "PUBLIC",
			SQL:      "SHOW FUNCTIONS",
			Enabled:  true,
			Schedule: "5 MINUTE",
		},

		ChildTask: &TaskSettings{
			Name:    childname,
			SQL:     "SELECT 1",
			Enabled: false,
			Comment: "initial state",
		},

		SoloTask: &TaskSettings{
			Name:    soloname,
			Schema:  "PUBLIC",
			SQL:     "SELECT 1",
			When:    "TRUE",
			Enabled: false,
		},
	}

	// Enables the Child and changes the SQL
	stepOne = &AccTaskTestSettings{
		WarehouseName: warehousename,
		DatabaseName:  databasename,

		RootTask: &TaskSettings{
			Name:     rootname,
			Schema:   "PUBLIC",
			SQL:      "SHOW FUNCTIONS",
			Enabled:  true,
			Schedule: "5 MINUTE",
		},

		ChildTask: &TaskSettings{
			Name:    childname,
			SQL:     "SELECT *",
			Enabled: true,
			Comment: "secondary state",
		},

		SoloTask: &TaskSettings{
			Name:    soloname,
			Schema:  "PUBLIC",
			SQL:     "SELECT *",
			When:    "TRUE",
			Enabled: true,
		},
	}

	// Changes Root Schedule and SQL
	stepTwo = &AccTaskTestSettings{
		WarehouseName: warehousename,
		DatabaseName:  databasename,

		RootTask: &TaskSettings{
			Name:     rootname,
			Schema:   "PUBLIC",
			SQL:      "SHOW TABLES",
			Enabled:  true,
			Schedule: "15 MINUTE",
		},

		ChildTask: &TaskSettings{
			Name:    childname,
			SQL:     "SELECT 1",
			Enabled: true,
			Comment: "third state",
		},

		SoloTask: &TaskSettings{
			Name:    soloname,
			Schema:  "PUBLIC",
			SQL:     "SELECT *",
			When:    "FALSE",
			Enabled: true,
		},
	}

	stepThree = &AccTaskTestSettings{
		WarehouseName: warehousename,
		DatabaseName:  databasename,

		RootTask: &TaskSettings{
			Name:     rootname,
			Schema:   "PUBLIC",
			SQL:      "SHOW FUNCTIONS",
			Enabled:  false,
			Schedule: "5 MINUTE",
		},

		ChildTask: &TaskSettings{
			Name:    childname,
			SQL:     "SELECT 1",
			Enabled: false,
			Comment: "reset",
		},

		SoloTask: &TaskSettings{
			Name:    soloname,
			Schema:  "PUBLIC",
			SQL:     "SELECT 1",
			When:    "TRUE",
			Enabled: true,
		},
	}
)

func Test_AccTask(t *testing.T) {
	resource.Test(t, resource.TestCase{
		Providers: providers(),
		Steps: []resource.TestStep{
			{
				Config: taskConfig(initialState),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.root_task", "enabled", true),
					checkBool("snowflake_task.child_task", "enabled", false),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", databasename),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", databasename),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", "PUBLIC"),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", "PUBLIC"),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql", initialState.RootTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql", initialState.ChildTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "after", initialState.RootTask.Name),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", initialState.ChildTask.Comment),
				),
			},
			{
				Config: taskConfig(stepOne),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.root_task", "enabled", true),
					checkBool("snowflake_task.child_task", "enabled", true),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", databasename),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", databasename),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", "PUBLIC"),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", "PUBLIC"),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql", stepOne.RootTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql", stepOne.ChildTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", stepOne.ChildTask.Comment),
				),
			},
			{
				Config: taskConfig(stepTwo),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.root_task", "enabled", true),
					checkBool("snowflake_task.child_task", "enabled", true),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", databasename),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", databasename),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", "PUBLIC"),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", "PUBLIC"),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql", stepTwo.RootTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql", stepTwo.ChildTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", stepTwo.ChildTask.Comment),
				),
			},
			{
				Config: taskConfig(stepThree),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.root_task", "enabled", false),
					checkBool("snowflake_task.child_task", "enabled", false),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", databasename),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", databasename),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", "PUBLIC"),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", "PUBLIC"),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql", stepThree.RootTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql", stepThree.ChildTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", stepThree.ChildTask.Comment),
				),
			},
			{
				Config: taskConfig(initialState),
				Check: resource.ComposeTestCheckFunc(
					checkBool("snowflake_task.root_task", "enabled", true),
					checkBool("snowflake_task.child_task", "enabled", false),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "name", rootname),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "name", childname),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "database", databasename),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "database", databasename),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "schema", "PUBLIC"),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "schema", "PUBLIC"),
					resource.TestCheckResourceAttr("snowflake_task.root_task", "sql", initialState.RootTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "sql", initialState.ChildTask.SQL),
					resource.TestCheckResourceAttr("snowflake_task.child_task", "comment", initialState.ChildTask.Comment),
				),
			},
		},
	})
}

func taskConfig(settings *AccTaskTestSettings) string {
	config, err := template.New("task_acceptance_test_config").Parse(`
resource "snowflake_warehouse" "test_wh" {
	name = "{{ .WarehouseName }}"
}

resource "snowflake_database" "test_db" {
	name = "{{ .DatabaseName }}"
}

resource "snowflake_task" "root_task" {
	name      = "{{ .RootTask.Name }}"
	database  = snowflake_database.test_db.name
	schema    = "{{ .RootTask.Schema }}"
	warehouse = snowflake_warehouse.test_wh.name
	sql       = "{{ .RootTask.SQL }}"
	enabled   = {{ .RootTask.Enabled }}
	schedule  = "{{ .RootTask.Schedule }}"
}

resource "snowflake_task" "child_task" {
	name      = "{{ .ChildTask.Name }}"
	database  = snowflake_task.root_task.database
	schema    = snowflake_task.root_task.schema
	warehouse = snowflake_task.root_task.warehouse
	sql       = "{{ .ChildTask.SQL }}"
	enabled   = {{ .ChildTask.Enabled }}
	after     = snowflake_task.root_task.name
	comment = "{{ .ChildTask.Comment }}"
}

resource "snowflake_task" "solo_task" {
	name      = "{{ .SoloTask.Name }}"
	database  = snowflake_database.test_db.name
	schema    = "{{ .SoloTask.Schema }}"
	warehouse = snowflake_warehouse.test_wh.name
	sql       = "{{ .SoloTask.SQL }}"
	enabled   = {{ .SoloTask.Enabled }}
	when      = "{{ .SoloTask.When }}"
}
	`)

	if err != nil {
		fmt.Println(err)
	}

	var result bytes.Buffer
	config.Execute(&result, settings)

	return result.String()
}

// type Inventory struct {
// 	Material string
// 	Count    uint
// }
// sweaters := Inventory{"wool", 17}
// tmpl, err := template.New("test").Parse("{{.Count}} items are made of {{.Material}}")
// if err != nil { panic(err) }
// err = tmpl.Execute(os.Stdout, sweaters)
// if err != nil { panic(err) }