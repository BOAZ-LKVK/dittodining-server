package sample

import (
	"errors"
	"github.com/BOAZ-LKVK/LKVK-server/domain/sample"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/apicontroller"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/customerrors"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/validate"
	sample_repository "github.com/BOAZ-LKVK/LKVK-server/repository/sample"
	"github.com/gofiber/fiber/v2"
)

// SampleAPIController implements apicontroller.APIController
var _ apicontroller.APIController = (*SampleAPIController)(nil)

type SampleAPIController struct {
	sampleRepository sample_repository.SampleRepository
}

func NewSampleAPIHandler(sampleRepository sample_repository.SampleRepository) *SampleAPIController {
	return &SampleAPIController{sampleRepository: sampleRepository}
}

func (c *SampleAPIController) Pattern() string {
	return "/samples"
}

func (c *SampleAPIController) Handlers() []*apicontroller.APIHandler {
	return []*apicontroller.APIHandler{
		apicontroller.NewAPIHandler("", fiber.MethodGet, c.listSamples()),
		apicontroller.NewAPIHandler("/:id", fiber.MethodGet, c.getSample()),
		apicontroller.NewAPIHandler("", fiber.MethodPost, c.createSample()),
		apicontroller.NewAPIHandler("/:id", fiber.MethodPut, c.updateSample()),
		apicontroller.NewAPIHandler("/:id", fiber.MethodDelete, c.deleteSample()),
	}
}

func (c *SampleAPIController) listSamples() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		samples, err := c.sampleRepository.FindAllSamples()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&ListSamplesResponse{Samples: samples})
	}
}

func (c *SampleAPIController) getSample() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sampleID := ctx.Params("id")
		sample, err := c.sampleRepository.FindOneSample(sampleID)
		if err != nil {
			if errors.Is(err, customerrors.ErrorSampleNotFound) {
				return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
			}

			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&GetSampleResponse{Sample: sample})
	}
}

func (c *SampleAPIController) deleteSample() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sampleID := ctx.Params("id")
		if err := c.sampleRepository.DeleteSample(sampleID); err != nil {
			if errors.Is(err, customerrors.ErrorSampleNotFound) {
				return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
			}

			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&DeleteSampleResponse{})
	}
}

func (c *SampleAPIController) createSample() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := new(CreateSampleRequest)
		if err := ctx.BodyParser(request); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		// 유효성 검사
		if err := validate.Validator.Struct(request); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		sample, err := c.sampleRepository.CreateSample(&sample.Sample{
			Name:  request.Name,
			Email: request.Email,
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&CreateSampleResponse{Sample: sample})
	}
}

func (c *SampleAPIController) updateSample() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sampleID := ctx.Params("id")

		request := new(UpdateSampleRequest)
		if err := ctx.BodyParser(request); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		// 유효성 검사
		if err := validate.Validator.Struct(request); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		sample, err := c.sampleRepository.UpdateSample(
			&sample.Sample{
				ID:    sampleID,
				Name:  request.Name,
				Email: request.Email,
			},
		)
		if err != nil {
			if errors.Is(err, customerrors.ErrorSampleNotFound) {
				return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
			}

			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&UpdateSampleResponse{
			Sample: sample,
		})
	}
}
