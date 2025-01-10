package ffmpeg

import (
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

func evaluateWildcards(request *ExecutionRequest) {
	request.Command = strings.ReplaceAll(request.Command, "${INPUT_FILE}", request.InputFile)
	request.Command = strings.ReplaceAll(request.Command, "${OUTPUT_FILE}", request.OutputFile)

	request.Command = strings.ReplaceAll(request.Command, "${DATE_YEAR}", time.Now().Format("2006"))
	request.Command = strings.ReplaceAll(request.Command, "${DATE_SHORTYEAR}", time.Now().Format("06"))
	request.Command = strings.ReplaceAll(request.Command, "${DATE_MONTH}", time.Now().Format("01"))
	request.Command = strings.ReplaceAll(request.Command, "${DATE_DAY}", time.Now().Format("02"))

	_, week := time.Now().ISOWeek()
	request.Command = strings.ReplaceAll(request.Command, "${DATE_WEEK}", strconv.Itoa(week))

	request.Command = strings.ReplaceAll(request.Command, "${TIME_HOUR}", time.Now().Format("15"))
	request.Command = strings.ReplaceAll(request.Command, "${TIME_MINUTE}", time.Now().Format("04"))
	request.Command = strings.ReplaceAll(request.Command, "${TIME_SECOND}", time.Now().Format("05"))

	request.Command = strings.ReplaceAll(request.Command, "${TIMESTAMP_SECONDS}", strconv.FormatInt(time.Now().Unix(), 10))
	request.Command = strings.ReplaceAll(request.Command, "${TIMESTAMP_MILLISECONDS}", strconv.FormatInt(time.Now().UnixMilli(), 10))
	request.Command = strings.ReplaceAll(request.Command, "${TIMESTAMP_MICROSECONDS}", strconv.FormatInt(time.Now().UnixMicro(), 10))
	request.Command = strings.ReplaceAll(request.Command, "${TIMESTAMP_NANOSECONDS}", strconv.FormatInt(time.Now().UnixNano(), 10))

	request.Command = strings.ReplaceAll(request.Command, "${OS_NAME}", runtime.GOOS)
	request.Command = strings.ReplaceAll(request.Command, "${OS_ARCH}", runtime.GOARCH)

	request.Command = strings.ReplaceAll(request.Command, "${UUID}", uuid.NewString())
}
