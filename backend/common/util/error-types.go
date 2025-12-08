package util

import (
	"errors"
	"fmt"
	"net/http"
)

func ErrRequiredInputMissing(field string) error {
	return errors.New("required input missing: " + field)
}

func ErrUnchangeable(field string) error {
	return errors.New("field unchangeable after creation: " + field)
}

func ErrHeaderMissing(field string) error {
	return errors.New("header missing: " + field)
}

func ErrAccessoryMappedToUser(phone string) error {
	return fmt.Errorf("accessory is already mapped to user with phone: %s", phone)
}

var (
	ErrExpiredToken          = errors.New("expired token")
	ErrInvalidToken          = errors.New("invalid token")
	ErrInvalidOtp            = errors.New("invalid otp")
	ErrInternal              = errors.New("internal error")
	ErrTokenMissing          = errors.New("token missing")
	ErrContextMissing        = errors.New("context missing")
	ErrUnauthorized          = errors.New("unauthorized")
	ErrTrainerAccessRequired = errors.New("trainer access required")
	ErrInvalidRole           = errors.New("invalid role")
	ErrEnvironment           = errors.New("environment error")
	ErrChatAlreadyExists     = errors.New("chat already exists")
	ErrChatNotFound          = errors.New("chat not found")
	ErrInvalidTimeFormat     = errors.New("invalid time format")

	ErrDatabase                              = errors.New("database error")
	ErrOdoLimit                              = errors.New("odo add limit reached")
	ErrEmptyResult                           = errors.New("empty result")
	ErrTooManyRequests                       = errors.New("too many requests")
	ErrEntryExists                           = errors.New("entry already exists")
	ErrInvalidUser                           = errors.New("invalid user")
	ErrSessionExpired                        = errors.New("session expired (something fishy may be happening)")
	ErrSmsUnableToSend                       = errors.New("unable to send sms")
	ErrExpiredOtp                            = errors.New("expired otp")
	ErrVinAlreadyPaired                      = errors.New("vin already paired")
	ErrVinNotPaired                          = errors.New("vin not paired to your profile")
	ErrInvalidPin                            = errors.New("invalid pin")
	ErrFileNotFound                          = errors.New("file not found")
	ErrEmailAlreadyVerified                  = errors.New("email already verified")
	ErrEmailNotVerified                      = errors.New("email not verified")
	ErrUrlParamsMissing                      = errors.New("url params missing")
	ErrModelDoesNotExist                     = errors.New("model does not exist")
	ErrInvalidFileFormat                     = errors.New("invalid file format")
	ErrSamePin                               = errors.New("same pin")
	ErrTooManyEntries                        = errors.New("too many entries")
	ErrInvalidObjectId                       = errors.New("invalid object id")
	ErrLocationNotFound                      = errors.New("location not found")
	ErrCouldNotSaveNotification              = errors.New("could not save notification")
	ErrCouldNotReestablishDatabaseConnection = errors.New("could not reestablish database connection")
	ErrEventCreationLimitExceeded            = errors.New("event creation limit exceeded")
	ErrAlreadyRegistered                     = errors.New("already registered")
	ErrNotValidRequest                       = errors.New("not valid request")
	ErrMultipartFileMissing                  = errors.New("multipart file missing")
	ErrNotInLocation                         = errors.New("not in location")
	ErrInvalidUuid                           = errors.New("invalid uuid")
	ErrGrpcConnFailed                        = errors.New("grpc connection failed")
	ErrGrpcQueryError                        = errors.New("grpc could not execute query")
	ErrNoGpsMapped                           = errors.New("no gps mapped to this profile")
	ErrGrpcCouldNotExecuteQuery              = errors.New("grpc could not execute query")
	ErrUserAlreadyExsists                    = errors.New("user already exsists")
	ErrUserDoesNotExist                      = errors.New("user does not exist")
	ErrWrongPassword                         = errors.New("wrong password")
	ErrNotEnoughCoins                        = errors.New("not enough coins for voucher")
	ErrInputInErrData                        = errors.New("errdata should contain either 0 or 1")
	ErrFramenumber                           = errors.New("invalid frame number input")
	ErrFrameExists                           = errors.New("frame number already exists")
	ErrInvalidQueryParams                    = errors.New("invalid query parameters")
	ErrTokenGenError                         = errors.New("error in generating token")
	ErrForbidden                             = errors.New("not allowed to access this resource")
	ErrAccessory                             = errors.New("accessory type does not exist")
	ErrAccessoryMapped                       = errors.New("accessory mapped with other user")
	ErrTripIdMissing                         = errors.New("no tripId passed")
	ErrMultipleLogin                         = errors.New("multiple login")
	ErrMinLatLng                             = errors.New("at least two positions are required to create a map")
	ErrMapUrl                                = errors.New("error fetching the image")
	ErrWebsocket                             = errors.New("error connecting to websocket")
	AccessoryAlreadyRequested                = errors.New("accessory already requested")
	ErrControllerIdAccessoryMapped           = errors.New("controllerId already exist with other framenumber")
	ErrBatteryIdAccessoryMapped              = errors.New("batteryId already exist with other framenumber")
	ErrMotorIdAccessoryMapped                = errors.New("motorId already exist with other framenumber")
	ErrNotCompressedImage                    = errors.New("image is not compressed")
	ErrImageNotUploaded                      = errors.New("image not uploaded")
	ErrAssessmentNotFound                    = errors.New("assessment not found")
	ErrAssessmentNotScheduled                = errors.New("assessment is not in scheduled status")
	ErrAssessmentNotCompleted                = errors.New("user has not completed assessment yet")
	ErrScheduleInPast                        = errors.New("scheduled time cannot be in the past")
	ErrSameScheduleTime                      = errors.New("new scheduled time is same as current time")
	ErrTrainerNotAvailable                   = errors.New("trainer is not available at the requested time")
	ErrNoTrainerAvailable                    = errors.New("no trainer is available at the requested time")
	ErrTrainerNotFound                       = errors.New("trainer not found")
	ErrClassNotFound                         = errors.New("class not found")
	ErrUserHasActiveAssessment               = errors.New("user already has an active assessment scheduled")
	ErrInvalidSessionType                    = errors.New("invalid session type for this assessment")
	ErrActiveAssessmentNotFound              = errors.New("no active assessment found for user")
	ErrAlreadyAssigned                       = errors.New("unauthorized: user is already assigned to a trainer")
	ErrNotAbleToAssignTrainer                = errors.New("not able to assign trainer")
	ErrInvalidFeedback                       = errors.New("invalid feedback provided")
	ErrNoMealLogCreated                      = errors.New("no meal log created for this user")
	ErrFailedToDecode                        = errors.New("failed to decode image, please upload a valid image format (png, jpeg, jpg)")
	ErrNoActiveMembership                    = errors.New("no active membership or error checking package eligibility")
	ErrNoMealLogAllowed                      = errors.New("your package does not allow meal logging")
	ErrFailedToFetchFeedback                 = errors.New("failed to fetch feedback for the meal log")
	ErrNoVideoFile                           = errors.New("video file is required")
	ErrNoExerciseVideoUploadAllowed          = errors.New("not eligible for exercise uploads")
	ErrUserHasActiveMembership               = errors.New("user already has an active membership")
	ErrItemAlreadyInCart                     = errors.New("item already in cart")
	ErrNutritionPlanUploadNotAllowed         = errors.New("nutrition plan upload is only allowed for SELF_TRAINING_EXERCISE_NUTRITION packages")
	ErrPlanNotReady                          = errors.New("plan is not ready for this membership")
	ErrNotEligibleForThisFeature             = errors.New("not eligible for this feature, please upgrade your membership")
	ErrAppointmentExistsForThisWeek          = errors.New("weekly appointment already exists for this week")
	ErrOverlappingAppointmentForTrainer      = errors.New("trainer already has an overlapping appointment at this time")
	ErrAppointmentNotFound                   = errors.New("appointment not found")
	ErrWeeklyAppointmentNotCreated           = errors.New("weekly appointment not created")
	ErrTimingError                           = errors.New("end time must be after start time")
	ErrAppointmentInFuture                   = errors.New("appointment must be in the future")
	ErrUnauthorizedForAppointment            = errors.New("unauthorized to reschedule this appointment")
	ErrTrainerNotAssigned                    = errors.New("trainer is not assigned to this membership")
	ErrInvalidRescheduledByValue             = errors.New("invalid 'rescheduled_by' value. Must be 'CLIENT' or 'TRAINER'")
)

var CustomErrorType = map[error]int{
	ErrExpiredToken:                          http.StatusUnauthorized,
	ErrInvalidToken:                          http.StatusUnauthorized,
	ErrInvalidOtp:                            419, // defining 419 for when a users token doesn't exist in cache
	ErrInternal:                              http.StatusInternalServerError,
	ErrTokenMissing:                          http.StatusUnauthorized,
	ErrContextMissing:                        http.StatusInternalServerError,
	ErrUnauthorized:                          http.StatusForbidden,
	ErrInvalidRole:                           http.StatusBadRequest,
	ErrEnvironment:                           http.StatusInternalServerError,
	ErrDatabase:                              http.StatusInternalServerError,
	ErrOdoLimit:                              http.StatusNotAcceptable,
	ErrEmptyResult:                           http.StatusBadRequest,
	ErrTooManyRequests:                       http.StatusTooManyRequests,
	ErrEntryExists:                           http.StatusBadRequest,
	ErrInvalidUser:                           http.StatusBadRequest,
	ErrSessionExpired:                        420, // defining 420 for when a users token doesn't exist in cache (someone else has created a new session)
	ErrSmsUnableToSend:                       430, // defining 430 for when an SMS is unable to send
	ErrExpiredOtp:                            429,
	ErrVinAlreadyPaired:                      427, // defining 427 for when a vin is already paired
	ErrVinNotPaired:                          http.StatusBadRequest,
	ErrFileNotFound:                          http.StatusInternalServerError,
	ErrEmailAlreadyVerified:                  http.StatusBadRequest,
	ErrInvalidPin:                            http.StatusBadRequest,
	ErrEmailNotVerified:                      http.StatusBadRequest,
	ErrModelDoesNotExist:                     http.StatusBadRequest,
	ErrChatAlreadyExists:                     http.StatusConflict,
	ErrChatNotFound:                          http.StatusNotFound,
	ErrInvalidTimeFormat:                     http.StatusBadRequest,
	ErrInvalidFileFormat:                     http.StatusBadRequest,
	ErrUrlParamsMissing:                      http.StatusBadRequest,
	ErrSamePin:                               http.StatusBadRequest,
	ErrTooManyEntries:                        http.StatusBadRequest,
	ErrInvalidObjectId:                       http.StatusBadRequest,
	ErrLocationNotFound:                      http.StatusBadRequest,
	ErrCouldNotSaveNotification:              http.StatusInternalServerError,
	ErrCouldNotReestablishDatabaseConnection: http.StatusInternalServerError,
	ErrInputInErrData:                        http.StatusNotAcceptable,
	ErrFramenumber:                           422,
	ErrFrameExists:                           427,
	ErrInvalidQueryParams:                    http.StatusBadRequest,
	ErrTokenGenError:                         http.StatusInternalServerError,
	ErrForbidden:                             http.StatusForbidden,
	ErrAccessory:                             http.StatusBadRequest,
	ErrAccessoryMapped:                       409,
	ErrTripIdMissing:                         http.StatusBadRequest,
	ErrMultipleLogin:                         http.StatusConflict,
	ErrMinLatLng:                             http.StatusBadRequest,
	ErrMapUrl:                                403,
	ErrWebsocket:                             http.StatusInternalServerError,
	ErrUserDoesNotExist:                      http.StatusForbidden,
	AccessoryAlreadyRequested:                409,
	ErrControllerIdAccessoryMapped:           400,
	ErrBatteryIdAccessoryMapped:              400,
	ErrMotorIdAccessoryMapped:                400,
	ErrTrainerAccessRequired:                 http.StatusForbidden,
	ErrAssessmentNotFound:                    http.StatusNotFound,
	ErrAssessmentNotScheduled:                http.StatusBadRequest,
	ErrAssessmentNotCompleted:                http.StatusBadRequest,
	ErrScheduleInPast:                        http.StatusBadRequest,
	ErrSameScheduleTime:                      http.StatusBadRequest,
	ErrTrainerNotAvailable:                   http.StatusConflict,
	ErrNoTrainerAvailable:                    http.StatusConflict,
	ErrTrainerNotFound:                       http.StatusNotFound,
	ErrUserHasActiveAssessment:               http.StatusConflict,
	ErrInvalidSessionType:                    http.StatusBadRequest,
	ErrActiveAssessmentNotFound:              http.StatusNotFound,
	ErrClassNotFound:                         http.StatusNotFound,
	ErrUserHasActiveMembership:               http.StatusConflict,
	ErrItemAlreadyInCart:                     http.StatusConflict,
	ErrNutritionPlanUploadNotAllowed:         http.StatusBadRequest,
	ErrPlanNotReady:                          http.StatusBadRequest,
	ErrTrainerNotAssigned:                    http.StatusBadRequest,
}
