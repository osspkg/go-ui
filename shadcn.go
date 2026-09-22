/*
 *  Copyright (c) 2026 Mikhail Knyazhev <markus621@yandex.com>. All rights reserved.
 *  Use of this source code is governed by a BSD 3-Clause license that can be found in the LICENSE file.
 */

package ui

// ShadcnComponent is a host-owned component name from the official shadcn/ui
// catalog. Register a matching React implementation before rendering it.
type ShadcnComponent string

const (
	ShadcnAccordion       ShadcnComponent = "accordion"
	ShadcnAlert           ShadcnComponent = "alert"
	ShadcnAlertDialog     ShadcnComponent = "alert-dialog"
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

func Accordion(id string) *NodeBuilder       { return Shadcn(ShadcnAccordion, id) }
func Alert(id string) *NodeBuilder           { return Shadcn(ShadcnAlert, id) }
func AlertDialog(id string) *NodeBuilder     { return Shadcn(ShadcnAlertDialog, id) }
func AspectRatio(id string) *NodeBuilder     { return Shadcn(ShadcnAspectRatio, id) }
func Attachment(id string) *NodeBuilder      { return Shadcn(ShadcnAttachment, id) }
func Avatar(id string) *NodeBuilder          { return Shadcn(ShadcnAvatar, id) }
func Badge(id string) *NodeBuilder           { return Shadcn(ShadcnBadge, id) }
func Breadcrumb(id string) *NodeBuilder      { return Shadcn(ShadcnBreadcrumb, id) }
func Bubble(id string) *NodeBuilder          { return Shadcn(ShadcnBubble, id) }
func ButtonGroup(id string) *NodeBuilder     { return Shadcn(ShadcnButtonGroup, id) }
func Calendar(id string) *NodeBuilder        { return Shadcn(ShadcnCalendar, id) }
func Carousel(id string) *NodeBuilder        { return Shadcn(ShadcnCarousel, id) }
func Chart(id string) *NodeBuilder           { return Shadcn(ShadcnChart, id) }
func Checkbox(id string) *NodeBuilder        { return Shadcn(ShadcnCheckbox, id) }
func Collapsible(id string) *NodeBuilder     { return Shadcn(ShadcnCollapsible, id) }
func Combobox(id string) *NodeBuilder        { return Shadcn(ShadcnCombobox, id) }
func Command(id string) *NodeBuilder         { return Shadcn(ShadcnCommand, id) }
func ContextMenu(id string) *NodeBuilder     { return Shadcn(ShadcnContextMenu, id) }
func DataTable(id string) *NodeBuilder       { return Shadcn(ShadcnDataTable, id) }
func DatePicker(id string) *NodeBuilder      { return Shadcn(ShadcnDatePicker, id) }
func Dialog(id string) *NodeBuilder          { return Shadcn(ShadcnDialog, id) }
func Direction(id string) *NodeBuilder       { return Shadcn(ShadcnDirection, id) }
func Drawer(id string) *NodeBuilder          { return Shadcn(ShadcnDrawer, id) }
func DropdownMenu(id string) *NodeBuilder    { return Shadcn(ShadcnDropdownMenu, id) }
func Empty(id string) *NodeBuilder           { return Shadcn(ShadcnEmpty, id) }
func Field(id string) *NodeBuilder           { return Shadcn(ShadcnField, id) }
func HoverCard(id string) *NodeBuilder       { return Shadcn(ShadcnHoverCard, id) }
func Input(id string) *NodeBuilder           { return Shadcn(ShadcnInput, id) }
func InputGroup(id string) *NodeBuilder      { return Shadcn(ShadcnInputGroup, id) }
func InputOTP(id string) *NodeBuilder        { return Shadcn(ShadcnInputOTP, id) }
func Item(id string) *NodeBuilder            { return Shadcn(ShadcnItem, id) }
func Kbd(id string) *NodeBuilder             { return Shadcn(ShadcnKbd, id) }
func Label(id string) *NodeBuilder           { return Shadcn(ShadcnLabel, id) }
func Marker(id string) *NodeBuilder          { return Shadcn(ShadcnMarker, id) }
func Menubar(id string) *NodeBuilder         { return Shadcn(ShadcnMenubar, id) }
func Message(id string) *NodeBuilder         { return Shadcn(ShadcnMessage, id) }
func MessageScroller(id string) *NodeBuilder { return Shadcn(ShadcnMessageScroller, id) }
func NativeSelect(id string) *NodeBuilder    { return Shadcn(ShadcnNativeSelect, id) }
func NavigationMenu(id string) *NodeBuilder  { return Shadcn(ShadcnNavigationMenu, id) }
func Pagination(id string) *NodeBuilder      { return Shadcn(ShadcnPagination, id) }
func Popover(id string) *NodeBuilder         { return Shadcn(ShadcnPopover, id) }
func Progress(id string) *NodeBuilder        { return Shadcn(ShadcnProgress, id) }
func Questionnaire(id string) *NodeBuilder   { return Shadcn(ShadcnQuestionnaire, id) }
func RadioGroup(id string) *NodeBuilder      { return Shadcn(ShadcnRadioGroup, id) }
func Resizable(id string) *NodeBuilder       { return Shadcn(ShadcnResizable, id) }
func ScrollArea(id string) *NodeBuilder      { return Shadcn(ShadcnScrollArea, id) }
func Select(id string) *NodeBuilder          { return Shadcn(ShadcnSelect, id) }
func Separator(id string) *NodeBuilder       { return Shadcn(ShadcnSeparator, id) }
func Sheet(id string) *NodeBuilder           { return Shadcn(ShadcnSheet, id) }
func Sidebar(id string) *NodeBuilder         { return Shadcn(ShadcnSidebar, id) }
func Skeleton(id string) *NodeBuilder        { return Shadcn(ShadcnSkeleton, id) }
func Slider(id string) *NodeBuilder          { return Shadcn(ShadcnSlider, id) }
func Spinner(id string) *NodeBuilder         { return Shadcn(ShadcnSpinner, id) }
func Switch(id string) *NodeBuilder          { return Shadcn(ShadcnSwitch, id) }
func Tabs(id string) *NodeBuilder            { return Shadcn(ShadcnTabs, id) }
func Textarea(id string) *NodeBuilder        { return Shadcn(ShadcnTextarea, id) }
func Toast(id string) *NodeBuilder           { return Shadcn(ShadcnToast, id) }
func Toggle(id string) *NodeBuilder          { return Shadcn(ShadcnToggle, id) }
func ToggleGroup(id string) *NodeBuilder     { return Shadcn(ShadcnToggleGroup, id) }
func Tooltip(id string) *NodeBuilder         { return Shadcn(ShadcnTooltip, id) }
func Typography(id string) *NodeBuilder      { return Shadcn(ShadcnTypography, id) }
