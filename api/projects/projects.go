package projects

import (
	database "github.com/ansible-semaphore/semaphore/db"
	"github.com/ansible-semaphore/semaphore/models"
	"github.com/gin-gonic/gin"
	"github.com/masterminds/squirrel"
)

func GetProjects(w http.ResponseWriter, r *http.Request) {
	user := context.Get(r, "user").(*models.User)

	query, args, _ := squirrel.Select("p.*").
		From("project as p").
		Join("project__user as pu on pu.project_id=p.id").
		Where("pu.user_id=?", user.ID).
		OrderBy("p.name").
		ToSql()

	var projects []models.Project
	if _, err := database.Mysql.Select(&projects, query, args...); err != nil {
		panic(err)
	}

	c.JSON(200, projects)
}

func AddProject(w http.ResponseWriter, r *http.Request) {
	var body models.Project
	user := context.Get(r, "user").(*models.User)

	if err := mulekick.Bind(w, r, &body); err != nil {
		return
	}

	err := body.CreateProject()
	if err != nil {
		panic(err)
	}

	if _, err := database.Mysql.Exec("insert into project__user set project_id=?, user_id=?, admin=1", body.ID, user.ID); err != nil {
		panic(err)
	}

	desc := "Project Created"
	if err := (models.Event{
		ProjectID:   &body.ID,
		Description: &desc,
	}.Insert()); err != nil {
		panic(err)
	}

	c.JSON(201, body)
}
