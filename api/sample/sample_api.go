package sample

import (
	"errors"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/apicontroller"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/customerrors"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/domain/sample"
	"github.com/BOAZ-LKVK/LKVK-server/pkg/validate"
	sample_repository "github.com/BOAZ-LKVK/LKVK-server/repository/sample"
	"github.com/gofiber/fiber/v2"
)

// SampleAPIHandler implements apicontroller.APIController
var _ apicontroller.APIController = (*SampleAPIHandler)(nil)

type SampleAPIHandler struct {
	sampleRepository sample_repository.SampleRepository
}

func NewSampleAPIHandler(sampleRepository sample_repository.SampleRepository) *SampleAPIHandler {
	return &SampleAPIHandler{sampleRepository: sampleRepository}
}

func (h *SampleAPIHandler) Pattern() string {
	return "/samples"
}

func (h *SampleAPIHandler) Handlers() []*apicontroller.APIHandler {
	return []*apicontroller.APIHandler{
		apicontroller.NewAPIHandler("", fiber.MethodGet, h.listSamples()),
		apicontroller.NewAPIHandler("/:id", fiber.MethodGet, h.getSample()),
		apicontroller.NewAPIHandler("", fiber.MethodPost, h.createSample()),
		apicontroller.NewAPIHandler("/:id", fiber.MethodPut, h.updateSample()),
		apicontroller.NewAPIHandler("/:id", fiber.MethodDelete, h.deleteSample()),
	}
}

func (h *SampleAPIHandler) listSamples() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		samples, err := h.sampleRepository.FindAllSamples()
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&ListSamplesResponse{Samples: samples})
	}
}

func (h *SampleAPIHandler) getSample() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sampleID := ctx.Params("id")
		sample, err := h.sampleRepository.FindOneSample(sampleID)
		if err != nil {
			if errors.Is(err, customerrors.ErrorSampleNotFound) {
				return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
			}

			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&GetSampleResponse{Sample: sample})
	}
}

func (h *SampleAPIHandler) deleteSample() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		sampleID := ctx.Params("id")
		if err := h.sampleRepository.DeleteSample(sampleID); err != nil {
			if errors.Is(err, customerrors.ErrorSampleNotFound) {
				return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
			}

			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&DeleteSampleResponse{})
	}
}

func (h *SampleAPIHandler) createSample() fiber.Handler {
	return func(ctx *fiber.Ctx) error {
		request := new(CreateSampleRequest)
		if err := ctx.BodyParser(request); err != nil {
			return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
		}

		// 유효성 검사
		if err := validate.Validator.Struct(request); err != nil {
			return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		sample, err := h.sampleRepository.CreateSample(&sample.Sample{
			Name:  request.Name,
			Email: request.Email,
		})
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		return ctx.JSON(&CreateSampleResponse{Sample: sample})
	}
}

func (h *SampleAPIHandler) updateSample() fiber.Handler {
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

		sample, err := h.sampleRepository.UpdateSample(
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
