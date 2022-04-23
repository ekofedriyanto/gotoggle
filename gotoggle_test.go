package gotoggle_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/ekofedriyanto/gotoggle"
	"github.com/stretchr/testify/assert"
	"go.chromium.org/luci/common/clock/testclock"
)

var now = time.Date(2006, 1, 2, 15, 4, 5, 6, time.UTC)
var flags = []gotoggle.Flag{
	{
		Flag: "GET /exists_with_no_time_data",
	},
	{
		Flag: "GET /exists_with_on_time_data_only_before_now",
		On:   time.Date(2006, 1, 2, 15, 4, 4, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_on_time_data_only_equal_now",
		On:   time.Date(2006, 1, 2, 15, 4, 5, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_on_time_data_only_after_now",
		On:   time.Date(2006, 1, 2, 15, 4, 6, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_off_time_data_only_before_now",
		Off:  time.Date(2006, 1, 2, 15, 4, 4, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_off_time_data_only_equal_now",
		Off:  time.Date(2006, 1, 2, 15, 4, 5, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_off_time_data_only_after_now",
		Off:  time.Date(2006, 1, 2, 15, 4, 6, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_off_and_on_time_data_equal",
		On:   time.Date(2006, 1, 2, 15, 4, 4, 6, time.UTC).Unix(),
		Off:  time.Date(2006, 1, 2, 15, 4, 4, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_off_and_on_time_data_before_now_and_on_greater_than_off",
		On:   time.Date(2006, 1, 2, 15, 4, 4, 6, time.UTC).Unix(),
		Off:  time.Date(2006, 1, 2, 15, 4, 3, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_off_and_on_time_data_after_now_and_on_greater_than_off",
		On:   time.Date(2006, 1, 2, 15, 4, 7, 6, time.UTC).Unix(),
		Off:  time.Date(2006, 1, 2, 15, 4, 6, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_off_and_on_time_data_between_now_and_on_greater_than_off",
		On:   time.Date(2006, 1, 2, 15, 4, 7, 6, time.UTC).Unix(),
		Off:  time.Date(2006, 1, 2, 15, 4, 4, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_off_and_on_time_data_off_is_now_and_on_greater_than_off",
		On:   time.Date(2006, 1, 2, 15, 4, 7, 6, time.UTC).Unix(),
		Off:  time.Date(2006, 1, 2, 15, 4, 5, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_off_and_on_time_data_on_is_now_and_on_greater_than_off",
		On:   time.Date(2006, 1, 2, 15, 4, 5, 6, time.UTC).Unix(),
		Off:  time.Date(2006, 1, 2, 15, 4, 4, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_off_and_on_time_data_before_now_and_on_less_than_off",
		On:   time.Date(2006, 1, 2, 15, 4, 3, 6, time.UTC).Unix(),
		Off:  time.Date(2006, 1, 2, 15, 4, 4, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_off_and_on_time_data_after_now_and_on_less_than_off",
		On:   time.Date(2006, 1, 2, 15, 4, 6, 6, time.UTC).Unix(),
		Off:  time.Date(2006, 1, 2, 15, 4, 7, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_off_and_on_time_data_between_now_and_on_less_than_off",
		On:   time.Date(2006, 1, 2, 15, 4, 4, 6, time.UTC).Unix(),
		Off:  time.Date(2006, 1, 2, 15, 4, 7, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_off_and_on_time_data_on_is_now_and_on_less_than_off",
		On:   time.Date(2006, 1, 2, 15, 4, 5, 6, time.UTC).Unix(),
		Off:  time.Date(2006, 1, 2, 15, 4, 7, 6, time.UTC).Unix(),
	},
	{
		Flag: "GET /exists_with_off_and_on_time_data_off_is_now_and_on_less_than_off",
		On:   time.Date(2006, 1, 2, 15, 4, 4, 6, time.UTC).Unix(),
		Off:  time.Date(2006, 1, 2, 15, 4, 5, 6, time.UTC).Unix(),
	},
}

type testData struct {
	inputString      string
	inputNow         time.Time
	isActiveExpected bool
}

func TestIsRouteActive(t *testing.T) {
	tests := map[string]testData{
		"Route Not Exists": {
			inputString:      "GET /not_exists",
			inputNow:         now,
			isActiveExpected: true,
		},
		"Route Exists With No Time Data": {
			inputString:      "GET /exists_with_no_time_data",
			inputNow:         now,
			isActiveExpected: true,
		},
		"Route Exists With On Time Data Only Before Now": {
			inputString:      "GET /exists_with_on_time_data_only_before_now",
			inputNow:         now,
			isActiveExpected: true,
		},
		"Route Exists With On Time Data Only Equal Now": {
			inputString:      "GET /exists_with_on_time_data_only_equal_now",
			inputNow:         now,
			isActiveExpected: true,
		},
		"Route Exists With On Time Data Only After Now": {
			inputString:      "GET /exists_with_on_time_data_only_after_now",
			inputNow:         now,
			isActiveExpected: false,
		},
		"Route Exists With Off Time Data Only Before Now": {
			inputString:      "GET /exists_with_off_time_data_only_before_now",
			inputNow:         now,
			isActiveExpected: false,
		},
		"Route Exists With Off Time Data Only Equal Now": {
			inputString:      "GET /exists_with_off_time_data_only_equal_now",
			inputNow:         now,
			isActiveExpected: false,
		},
		"Route Exists With Off Time Data Only After Now": {
			inputString:      "GET /exists_with_off_time_data_only_after_now",
			inputNow:         now,
			isActiveExpected: true,
		},
		"Route Exists With Off And On Time Data Equal": {
			inputString:      "GET /exists_with_off_and_on_time_data_equal",
			inputNow:         now,
			isActiveExpected: true,
		},
		"Route Exists With Off And On Time Data Before Now And On Greater Than Off": {
			inputString:      "GET /exists_with_off_and_on_time_data_before_now_and_on_greater_than_off",
			inputNow:         now,
			isActiveExpected: true,
		},
		"Route Exists With Off And On Time Data After Now And On Greater Than Off": {
			inputString:      "GET /exists_with_off_and_on_time_data_after_now_and_on_greater_than_off",
			inputNow:         now,
			isActiveExpected: false,
		},
		"Route Exists With Off And On Time Data Between Now And On Greater Than Off": {
			inputString:      "GET /exists_with_off_and_on_time_data_between_now_and_on_greater_than_off",
			inputNow:         now,
			isActiveExpected: false,
		},
		"Route Exists With Off And On Time Data Off Is Now And On Greater Than Off": {
			inputString:      "GET /exists_with_off_and_on_time_data_off_is_now_and_on_greater_than_off",
			inputNow:         now,
			isActiveExpected: false,
		},
		"Route Exists With Off And On Time Data On Is Now And On Greater Than Off": {
			inputString:      "GET /exists_with_off_and_on_time_data_on_is_now_and_on_greater_than_off",
			inputNow:         now,
			isActiveExpected: true,
		},
		"Route Exists With Off And On Time Data Before Now And On Less Than Off": {
			inputString:      "GET /exists_with_off_and_on_time_data_before_now_and_on_less_than_off",
			inputNow:         now,
			isActiveExpected: false,
		},
		"Route Exists With Off And On Time Data After Now And On Less Than Off": {
			inputString:      "GET /exists_with_off_and_on_time_data_after_now_and_on_less_than_off",
			inputNow:         now,
			isActiveExpected: false,
		},
		"Route Exists With Off And On Time Data Between Now And On Less Than Off": {
			inputString:      "GET /exists_with_off_and_on_time_data_between_now_and_on_less_than_off",
			inputNow:         now,
			isActiveExpected: true,
		},
		"Route Exists With Off And On Time Data On Is Now And On Less Than Off": {
			inputString:      "GET /exists_with_off_and_on_time_data_on_is_now_and_on_less_than_offs",
			inputNow:         now,
			isActiveExpected: true,
		},
		"Route Exists With Off And On Time Data Off Is Now And On Less Than Off": {
			inputString:      "GET /exists_with_off_and_on_time_data_off_is_now_and_on_less_than_off",
			inputNow:         now,
			isActiveExpected: false,
		},
	}

	toggleFlags := gotoggle.NewToggles(flags...)

	for name, testD := range tests {
		t.Run(name, func(tt *testing.T) {
			ctx, _ := testclock.UseTime(context.Background(), testD.inputNow)
			returnData := toggleFlags.IsActive(ctx, testD.inputString)
			assert.Equal(
				tt,
				testD.isActiveExpected,
				returnData,
				fmt.Sprintf(
					"IsActive not match (expected %v got %v)",
					testD.isActiveExpected,
					returnData,
				),
			)
		})
	}
}
