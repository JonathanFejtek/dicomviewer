package http

import (
	"dicomviewer/dicom"
	"errors"
	"image/png"
	"log/slog"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"

	dicomutil "github.com/suyashkumar/dicom"
	"github.com/suyashkumar/dicom/pkg/tag"
)

const tagQueryParamKey = "tag"

// dicomRecords contains a set of http handlers for managing DICOM files
type dicomFiles struct {
	fileRepository dicom.FileRepository
}

// Get an http handler to retrieve a raw DICOM file
func (f *dicomFiles) GetAll(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fileIDs, err := f.fileRepository.GetAll()
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	type Response struct {
		FileIDs []string `json:"fileIds"`
	}

	writeJSONResponse(w, Response{
		FileIDs: fileIDs,
	})
}

// Get an http handler to retrieve a raw DICOM file
func (f *dicomFiles) GetByID(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fileID, err := parseURLParam(r, "id")
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	file, httpStatus, err := f.getFile(fileID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		writeJSONError(w, httpStatus, err)
		return
	}

	http.ServeContent(w, r, file.ID, time.Now(), file.Raw())
}

// GetAsPNG an http handler to retreive a DICOM file as a grayscale PNG
func (f *dicomFiles) GetAsPNG(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	fileID, err := parseURLParam(r, "id")
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	shouldRemap, err := strconv.ParseBool(
		r.URL.Query().Get("remap"),
	)
	if err != nil {
		shouldRemap = true
	}

	file, httpStatus, err := f.getFile(fileID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		writeJSONError(w, httpStatus, err)
		return
	}

	image, err := file.PNG(dicom.PNGRemapPixels(shouldRemap))
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}

	err = png.Encode(w, image)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}
}

// SearchAttributes an http handler to search the attributes/elements of a
// DICOM file
func (f *dicomFiles) SearchAttributes(
	w http.ResponseWriter,
	r *http.Request,
) {
	ctx := r.Context()

	fileID, err := parseURLParam(r, "id")
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	dicomTags, err := f.parseTagQuery(r.URL.Query())
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}

	file, httpStatus, err := f.getFile(fileID)
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		writeJSONError(w, httpStatus, err)
		return
	}

	var elementsByTag map[string]*dicomutil.Element
	if len(dicomTags) > 0 {
		elementsByTag, err = file.FindElements(dicomTags...)
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
	} else {
		elementsByTag, err = file.AllElements()
		if err != nil {
			slog.ErrorContext(ctx, err.Error())
			writeJSONError(w, http.StatusInternalServerError, err)
			return
		}
	}

	type Response struct {
		ElementsByTag map[string]*dicomutil.Element `json:"elementsByTag"`
	}

	writeJSONResponse(
		w,
		Response{
			ElementsByTag: elementsByTag,
		},
	)
}

// Create an http handler to create a DICOM file
func (f *dicomFiles) Create(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	file, header, err := r.FormFile("file")
	if err != nil {
		slog.ErrorContext(ctx, err.Error())
		writeJSONError(w, http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	fileID := uuid.NewString()
	newFile := dicom.NewFile(
		fileID,
		header.Size,
		file,
	)

	if err := f.fileRepository.Create(
		newFile,
	); err != nil {
		slog.ErrorContext(ctx, err.Error())
		writeJSONError(w, http.StatusInternalServerError, err)
		return
	}

	type Response struct {
		FileID string `json:"fileId"`
	}

	writeJSONResponse(
		w,
		Response{
			FileID: fileID,
		},
	)
}

func (f *dicomFiles) getFile(id string) (*dicom.File, int, error) { // http status, err
	file, err := f.fileRepository.Get(id)
	if err != nil {
		status := http.StatusInternalServerError
		if errors.Is(err, dicom.ErrFileNotFound) {
			status = http.StatusNotFound
		}
		return nil, status, err
	}

	return file, http.StatusFound, nil
}

// parseTagQuery utility to parse url tag queries into tag objects
func (f dicomFiles) parseTagQuery(query url.Values) ([]tag.Tag, error) {

	queryStrings := query[tagQueryParamKey]

	var tags []tag.Tag
	for _, queryString := range queryStrings {
		tag, err := parseTag(queryString)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, nil
}

// parseTag adapted from helper in github.com/suyashkumar/dicom
func parseTag(tagString string) (tag.Tag, error) {
	parts := strings.Split(strings.Trim(tagString, "()"), ",")
	if len(parts) < 2 {
		return tag.Tag{}, errors.New("malformed tag")
	}
	group, err := strconv.ParseInt(parts[0], 16, 0)
	if err != nil {
		return tag.Tag{}, err
	}
	elem, err := strconv.ParseInt(parts[1], 16, 0)
	if err != nil {
		return tag.Tag{}, err
	}
	return tag.Tag{Group: uint16(group), Element: uint16(elem)}, nil
}
