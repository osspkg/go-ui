/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ui

// ShadcnComponent is a host-owned component name from the official shadcn/ui
// catalog. Register a matching React implementation before rendering it.
type ShadcnComponent string

// Shadcn constants identify host-owned shadcn/ui components.
const (
	// ShadcnAccordion through ShadcnTypography enumerate supported shadcn/ui components.
	ShadcnAccordion ShadcnComponent = "accordion"
	// ShadcnAlert identifies the alert component.
	ShadcnAlert ShadcnComponent = "alert"
	// ShadcnAlertDialog identifies the alert dialog component.
	ShadcnAlertDialog ShadcnComponent = "alert-dialog"
	// ShadcnAspectRatio identifies the aspect ratio component.
	ShadcnAspectRatio     ShadcnComponent = "aspect-ratio"
	ShadcnAttachment      ShadcnComponent = "attachment"
	ShadcnAvatar          ShadcnComponent = "avatar"
	ShadcnBadge           ShadcnComponent = "badge"
	ShadcnBreadcrumb      ShadcnComponent = "breadcrumb"
	ShadcnBubble          ShadcnComponent = "bubble"
	ShadcnButton          ShadcnComponent = "button"
	ShadcnButtonGroup     ShadcnComponent = "button-group"
	ShadcnCalendar        ShadcnComponent = "calendar"
	ShadcnCard            ShadcnComponent = "card"
	ShadcnCarousel        ShadcnComponent = "carousel"
	ShadcnChart           ShadcnComponent = "chart"
	ShadcnCheckbox        ShadcnComponent = "checkbox"
	ShadcnCollapsible     ShadcnComponent = "collapsible"
	ShadcnCombobox        ShadcnComponent = "combobox"
	ShadcnCommand         ShadcnComponent = "command"
	ShadcnContextMenu     ShadcnComponent = "context-menu"
	ShadcnDataTable       ShadcnComponent = "data-table"
	ShadcnDatePicker      ShadcnComponent = "date-picker"
	ShadcnDialog          ShadcnComponent = "dialog"
	ShadcnDirection       ShadcnComponent = "direction"
	ShadcnDrawer          ShadcnComponent = "drawer"
	ShadcnDropdownMenu    ShadcnComponent = "dropdown-menu"
	ShadcnEmpty           ShadcnComponent = "empty"
	ShadcnField           ShadcnComponent = "field"
	ShadcnHoverCard       ShadcnComponent = "hover-card"
	ShadcnInput           ShadcnComponent = "input"
	ShadcnInputGroup      ShadcnComponent = "input-group"
	ShadcnInputOTP        ShadcnComponent = "input-otp"
	ShadcnItem            ShadcnComponent = "item"
	ShadcnKbd             ShadcnComponent = "kbd"
	ShadcnLabel           ShadcnComponent = "label"
	ShadcnMarker          ShadcnComponent = "marker"
	ShadcnMenubar         ShadcnComponent = "menubar"
	ShadcnMessage         ShadcnComponent = "message"
	ShadcnMessageScroller ShadcnComponent = "message-scroller"
	ShadcnNativeSelect    ShadcnComponent = "native-select"
	ShadcnNavigationMenu  ShadcnComponent = "navigation-menu"
	ShadcnPagination      ShadcnComponent = "pagination"
	ShadcnPopover         ShadcnComponent = "popover"
	ShadcnProgress        ShadcnComponent = "progress"
	ShadcnQuestionnaire   ShadcnComponent = "questionnaire"
	ShadcnRadioGroup      ShadcnComponent = "radio-group"
	ShadcnResizable       ShadcnComponent = "resizable"
	ShadcnScrollArea      ShadcnComponent = "scroll-area"
	ShadcnSelect          ShadcnComponent = "select"
	ShadcnSeparator       ShadcnComponent = "separator"
	ShadcnSheet           ShadcnComponent = "sheet"
	ShadcnSidebar         ShadcnComponent = "sidebar"
	ShadcnSkeleton        ShadcnComponent = "skeleton"
	ShadcnSlider          ShadcnComponent = "slider"
	ShadcnSpinner         ShadcnComponent = "spinner"
	ShadcnSwitch          ShadcnComponent = "switch"
	ShadcnTable           ShadcnComponent = "table"
	ShadcnTabs            ShadcnComponent = "tabs"
	ShadcnTextarea        ShadcnComponent = "textarea"
	ShadcnToast           ShadcnComponent = "toast"
	ShadcnToggle          ShadcnComponent = "toggle"
	ShadcnToggleGroup     ShadcnComponent = "toggle-group"
	ShadcnTooltip         ShadcnComponent = "tooltip"
	ShadcnTypography      ShadcnComponent = "typography"
)

// Shadcn creates a node for a component from the shadcn/ui catalog.
func Shadcn(component ShadcnComponent, id string) *NodeBuilder { return node(string(component), id) }

// Accordion creates an accordion node.
func Accordion(id string) *NodeBuilder { return Shadcn(ShadcnAccordion, id) }

// Alert creates an alert node.
func Alert(id string) *NodeBuilder { return Shadcn(ShadcnAlert, id) }

// AlertDialog creates an alert dialog node.
func AlertDialog(id string) *NodeBuilder { return Shadcn(ShadcnAlertDialog, id) }

// AspectRatio creates an aspect ratio node.
func AspectRatio(id string) *NodeBuilder { return Shadcn(ShadcnAspectRatio, id) }

// Attachment creates an attachment node.
func Attachment(id string) *NodeBuilder { return Shadcn(ShadcnAttachment, id) }

// Avatar creates an avatar node.
func Avatar(id string) *NodeBuilder { return Shadcn(ShadcnAvatar, id) }

// Badge creates a badge node.
func Badge(id string) *NodeBuilder { return Shadcn(ShadcnBadge, id) }

// Breadcrumb creates a breadcrumb node.
func Breadcrumb(id string) *NodeBuilder { return Shadcn(ShadcnBreadcrumb, id) }

// Bubble creates a bubble node.
func Bubble(id string) *NodeBuilder { return Shadcn(ShadcnBubble, id) }

// ButtonGroup creates a button group node.
func ButtonGroup(id string) *NodeBuilder { return Shadcn(ShadcnButtonGroup, id) }

// Calendar creates a calendar node.
func Calendar(id string) *NodeBuilder { return Shadcn(ShadcnCalendar, id) }

// Carousel creates a carousel node.
func Carousel(id string) *NodeBuilder { return Shadcn(ShadcnCarousel, id) }

// Chart creates a chart node.
func Chart(id string) *NodeBuilder { return Shadcn(ShadcnChart, id) }

// Checkbox creates a checkbox node.
func Checkbox(id string) *NodeBuilder { return Shadcn(ShadcnCheckbox, id) }

// Collapsible creates a collapsible node.
func Collapsible(id string) *NodeBuilder { return Shadcn(ShadcnCollapsible, id) }

// Combobox creates a combobox node.
func Combobox(id string) *NodeBuilder { return Shadcn(ShadcnCombobox, id) }

// Command creates a command node.
func Command(id string) *NodeBuilder { return Shadcn(ShadcnCommand, id) }

// ContextMenu creates a context menu node.
func ContextMenu(id string) *NodeBuilder { return Shadcn(ShadcnContextMenu, id) }

// DataTable creates a data table node.
func DataTable(id string) *NodeBuilder { return Shadcn(ShadcnDataTable, id) }

// DatePicker creates a date picker node.
func DatePicker(id string) *NodeBuilder { return Shadcn(ShadcnDatePicker, id) }

// Dialog creates a dialog node.
func Dialog(id string) *NodeBuilder { return Shadcn(ShadcnDialog, id) }

// Direction creates a direction node.
func Direction(id string) *NodeBuilder { return Shadcn(ShadcnDirection, id) }

// Drawer creates a drawer node.
func Drawer(id string) *NodeBuilder { return Shadcn(ShadcnDrawer, id) }

// DropdownMenu creates a dropdown menu node.
func DropdownMenu(id string) *NodeBuilder { return Shadcn(ShadcnDropdownMenu, id) }

// Empty creates an empty-state node.
func Empty(id string) *NodeBuilder { return Shadcn(ShadcnEmpty, id) }

// Field creates a field node.
func Field(id string) *NodeBuilder { return Shadcn(ShadcnField, id) }

// HoverCard creates a hover card node.
func HoverCard(id string) *NodeBuilder { return Shadcn(ShadcnHoverCard, id) }

// Input creates an input node.
func Input(id string) *NodeBuilder { return Shadcn(ShadcnInput, id) }

// InputGroup creates an input group node.
func InputGroup(id string) *NodeBuilder { return Shadcn(ShadcnInputGroup, id) }

// InputOTP creates an input OTP node.
func InputOTP(id string) *NodeBuilder { return Shadcn(ShadcnInputOTP, id) }

// Item creates an item node.
func Item(id string) *NodeBuilder { return Shadcn(ShadcnItem, id) }

// Kbd creates a keyboard shortcut node.
func Kbd(id string) *NodeBuilder { return Shadcn(ShadcnKbd, id) }

// Label creates a label node.
func Label(id string) *NodeBuilder { return Shadcn(ShadcnLabel, id) }

// Marker creates a marker node.
func Marker(id string) *NodeBuilder { return Shadcn(ShadcnMarker, id) }

// Menubar creates a menubar node.
func Menubar(id string) *NodeBuilder { return Shadcn(ShadcnMenubar, id) }

// Message creates a message node.
func Message(id string) *NodeBuilder { return Shadcn(ShadcnMessage, id) }

// MessageScroller creates a message scroller node.
func MessageScroller(id string) *NodeBuilder { return Shadcn(ShadcnMessageScroller, id) }

// NativeSelect creates a native select node.
func NativeSelect(id string) *NodeBuilder { return Shadcn(ShadcnNativeSelect, id) }

// NavigationMenu creates a navigation menu node.
func NavigationMenu(id string) *NodeBuilder { return Shadcn(ShadcnNavigationMenu, id) }

// Pagination creates a pagination node.
func Pagination(id string) *NodeBuilder { return Shadcn(ShadcnPagination, id) }

// Popover creates a popover node.
func Popover(id string) *NodeBuilder { return Shadcn(ShadcnPopover, id) }

// Progress creates a progress node.
func Progress(id string) *NodeBuilder { return Shadcn(ShadcnProgress, id) }

// Questionnaire creates a questionnaire node.
func Questionnaire(id string) *NodeBuilder { return Shadcn(ShadcnQuestionnaire, id) }

// RadioGroup creates a radio group node.
func RadioGroup(id string) *NodeBuilder { return Shadcn(ShadcnRadioGroup, id) }

// Resizable creates a resizable node.
func Resizable(id string) *NodeBuilder { return Shadcn(ShadcnResizable, id) }

// ScrollArea creates a scroll area node.
func ScrollArea(id string) *NodeBuilder { return Shadcn(ShadcnScrollArea, id) }

// Select creates a select node.
func Select(id string) *NodeBuilder { return Shadcn(ShadcnSelect, id) }

// Separator creates a separator node.
func Separator(id string) *NodeBuilder { return Shadcn(ShadcnSeparator, id) }

// Sheet creates a sheet node.
func Sheet(id string) *NodeBuilder { return Shadcn(ShadcnSheet, id) }

// Sidebar creates a sidebar node.
func Sidebar(id string) *NodeBuilder { return Shadcn(ShadcnSidebar, id) }

// Skeleton creates a skeleton node.
func Skeleton(id string) *NodeBuilder { return Shadcn(ShadcnSkeleton, id) }

// Slider creates a slider node.
func Slider(id string) *NodeBuilder { return Shadcn(ShadcnSlider, id) }

// Spinner creates a spinner node.
func Spinner(id string) *NodeBuilder { return Shadcn(ShadcnSpinner, id) }

// Switch creates a switch node.
func Switch(id string) *NodeBuilder { return Shadcn(ShadcnSwitch, id) }

// Tabs creates a tabs node.
func Tabs(id string) *NodeBuilder { return Shadcn(ShadcnTabs, id) }

// Textarea creates a textarea node.
func Textarea(id string) *NodeBuilder { return Shadcn(ShadcnTextarea, id) }

// Toast creates a toast node.
func Toast(id string) *NodeBuilder { return Shadcn(ShadcnToast, id) }

// Toggle creates a toggle node.
func Toggle(id string) *NodeBuilder { return Shadcn(ShadcnToggle, id) }

// ToggleGroup creates a toggle group node.
func ToggleGroup(id string) *NodeBuilder { return Shadcn(ShadcnToggleGroup, id) }

// Tooltip creates a tooltip node.
func Tooltip(id string) *NodeBuilder { return Shadcn(ShadcnTooltip, id) }

// Typography creates a typography node.
func Typography(id string) *NodeBuilder { return Shadcn(ShadcnTypography, id) }
