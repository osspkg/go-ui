package ui

import "testing"

func TestUnit_ShadcnBuildersProduceCatalogNames(t *testing.T) {
	tests := []struct {
		name      string
		component ShadcnComponent
	}{
		{"accordion", ShadcnAccordion},
		{"alert", ShadcnAlert},
		{"alert dialog", ShadcnAlertDialog},
		{"aspect ratio", ShadcnAspectRatio},
		{"attachment", ShadcnAttachment},
		{"avatar", ShadcnAvatar},
		{"badge", ShadcnBadge},
		{"breadcrumb", ShadcnBreadcrumb},
		{"bubble", ShadcnBubble},
		{"button", ShadcnButton},
		{"button group", ShadcnButtonGroup},
		{"calendar", ShadcnCalendar},
		{"card", ShadcnCard},
		{"carousel", ShadcnCarousel},
		{"chart", ShadcnChart},
		{"checkbox", ShadcnCheckbox},
		{"collapsible", ShadcnCollapsible},
		{"combobox", ShadcnCombobox},
		{"command", ShadcnCommand},
		{"context menu", ShadcnContextMenu},
		{"data table", ShadcnDataTable},
		{"date picker", ShadcnDatePicker},
		{"dialog", ShadcnDialog},
		{"direction", ShadcnDirection},
		{"drawer", ShadcnDrawer},
		{"dropdown menu", ShadcnDropdownMenu},
		{"empty", ShadcnEmpty},
		{"field", ShadcnField},
		{"hover card", ShadcnHoverCard},
		{"input", ShadcnInput},
		{"input group", ShadcnInputGroup},
		{"input otp", ShadcnInputOTP},
		{"item", ShadcnItem},
		{"kbd", ShadcnKbd},
		{"label", ShadcnLabel},
		{"marker", ShadcnMarker},
		{"menubar", ShadcnMenubar},
		{"message", ShadcnMessage},
		{"message scroller", ShadcnMessageScroller},
		{"native select", ShadcnNativeSelect},
		{"navigation menu", ShadcnNavigationMenu},
		{"pagination", ShadcnPagination},
		{"popover", ShadcnPopover},
		{"progress", ShadcnProgress},
		{"questionnaire", ShadcnQuestionnaire},
		{"radio group", ShadcnRadioGroup},
		{"resizable", ShadcnResizable},
		{"scroll area", ShadcnScrollArea},
		{"select", ShadcnSelect},
		{"separator", ShadcnSeparator},
		{"sheet", ShadcnSheet},
		{"sidebar", ShadcnSidebar},
		{"skeleton", ShadcnSkeleton},
		{"slider", ShadcnSlider},
		{"spinner", ShadcnSpinner},
		{"switch", ShadcnSwitch},
		{"table", ShadcnTable},
		{"tabs", ShadcnTabs},
		{"textarea", ShadcnTextarea},
		{"toast", ShadcnToast},
		{"toggle", ShadcnToggle},
		{"toggle group", ShadcnToggleGroup},
		{"tooltip", ShadcnTooltip},
		{"typography", ShadcnTypography},
	}
	for _, test := range tests {
		t.Run(
			test.name, func(t *testing.T) {
				if got := Shadcn(test.component, "node").Build().Component; got != string(test.component) {
					t.Fatalf("component = %q, want %q", got, test.component)
				}
			},
		)
	}
}

func TestUnit_ShadcnConvenienceBuilder(t *testing.T) {
	if got := AlertDialog("confirm").Build().Component; got != string(ShadcnAlertDialog) {
		t.Fatalf("AlertDialog() component = %q, want %q", got, ShadcnAlertDialog)
	}
}
