package attendance

import (
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
	"gitlab.ecapture.com.co/gitlab-instance/gitlab-instance-cea63b52/e-capture/indra/api-indra-admin/internal/logger"
	"gitlab.ecapture.com.co/gitlab-instance/gitlab-instance-cea63b52/e-capture/indra/api-indra-admin/internal/middleware"
	"gitlab.ecapture.com.co/gitlab-instance/gitlab-instance-cea63b52/e-capture/indra/api-indra-admin/internal/models"
	"gitlab.ecapture.com.co/gitlab-instance/gitlab-instance-cea63b52/e-capture/indra/api-indra-admin/internal/msgs"
	"gitlab.ecapture.com.co/gitlab-instance/gitlab-instance-cea63b52/e-capture/indra/api-indra-admin/pkg/auth"
	"gitlab.ecapture.com.co/gitlab-instance/gitlab-instance-cea63b52/e-capture/indra/api-indra-admin/pkg/entity"
	"net/http"
)

type handlerAttendance struct {
	dB   *sqlx.DB
	user *models.User
	txID string
}

func (h *handlerAttendance) CreateAttendance(c *fiber.Ctx) error {
	msg := msgs.Model{}
	res := ResponseAttendance{Error: true}
	srvEntity := entity.NewServerEntity(h.dB, h.user, h.txID)
	srvAuth := auth.NewServerAuth(h.dB, h.user, h.txID)
	rqAAttendance := RequestAttendance{}

	err := c.BodyParser(&rqAAttendance)
	if err != nil {
		logger.Error.Printf("couldn't bind model BodyParser: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(1)
		res.Msg = "couldn't bind model RequestMetadata"
		return c.Status(http.StatusAccepted).JSON(res)
	}

	usr, err := middleware.GetUser(c)
	if err != nil {
		res.Error = true
		res.Code = 99
		res.Msg = "Error in token"
		return c.Status(http.StatusOK).JSON(res)
	}

	usrUser, cod, err := srvAuth.Users.GetUserByCodeStudent(rqAAttendance.CodeStudent)
	if err != nil {
		logger.Error.Printf("Couldn't insert suffragers: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(cod)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	req, cod, err := srvEntity.Attendance.CreateAttendance(usrUser.ID, rqAAttendance.IdEvent, 0, 0, usr.ID)
	if err != nil {
		logger.Error.Printf("Couldn't insert CreateAttendance: %v", err)
		res.Code, res.Type, res.Msg = msg.GetByCode(cod)
		return c.Status(http.StatusAccepted).JSON(res)
	}

	res.Data = req
	res.Error = false

	return c.Status(http.StatusOK).JSON(res)
}
