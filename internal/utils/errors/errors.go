package errors

import "GoBackendExploreMovieTracker/internal/utils"

var (
	ERROR_STATUS_BAD_REQUEST           = utils.Envelope{"error": "bad request payload"}
	ERROR_STATUS_INTERNAL_SERVER_ERROR = utils.Envelope{"error": "internal server error"}
	ERROR_STATUS_NOT_FOUND             = utils.Envelope{"error": "resource not found"}
)
