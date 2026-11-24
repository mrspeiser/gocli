# CLI Styling and UI Enhancements

This document describes the visual and user experience improvements made to the GOCLI ETL Tool.

## Overview

The CLI has been enhanced with modern terminal styling including:
- 🎨 **Color-coded output** for better readability
- 📊 **Progress bars** for HTTP operations
- 🎯 **Icons** for visual context
- 📋 **Formatted tables** for data display
- ✨ **Styled headers and sections**

## Features

### 1. Progress Bars

Progress bars are now displayed during:
- **Downloads** (HTTP GET operations from sources)
- **Uploads** (HTTP POST operations to destinations)

The progress bars show:
- Current progress with visual bar
- Bytes transferred
- Real-time updates
- Completion status

### 2. Color Scheme

The application uses a consistent color scheme:
- **Green** - Success messages and confirmations
- **Red** - Error messages
- **Yellow** - Warnings and required actions
- **Cyan** - Informational messages
- **Blue** - Prompts and user input
- **Magenta** - Headers and titles

### 3. Icons

Visual icons enhance the user experience:
- 🔄 **Pipeline** operations
- 📥 **Source** (download/extract)
- 📤 **Destination** (upload/load)
- ⚙️ **Transform** operations
- 📋 **Sequence** management
- 💾 **State** persistence
- ✅ **Success** indicators
- ❌ **Error** indicators
- ➕ **Add/Create** actions
- 🗑️ **Delete** actions
- 👁️ **View/Show** actions
- 🚀 **Run/Execute** actions

### 4. Formatted Tables

Data is displayed in well-formatted tables with:
- Header rows with column names
- Color-coded cells
- Proper alignment
- Visual separation

### 5. Styled Output

#### Application Header
The application starts with an ASCII art banner displaying "GO CLI" in large letters.

#### Menu Display
Menus are presented in table format with:
- Option numbers/letters
- Descriptions with relevant icons
- Clear visual hierarchy

#### Status Messages
All status messages are styled:
- Success messages with green checkmarks
- Errors with red indicators
- Warnings with yellow highlights
- Info messages with blue styling

### 6. User Prompts

Interactive prompts are styled with:
- Blue arrow (→) indicator
- Clear label text
- Consistent formatting

## Libraries Used

The styling is powered by:
- **pterm** - Beautiful terminal output with colors, tables, and formatting
- **fatih/color** - Additional color support
- **progressbar** - HTTP operation progress bars

## Example Output

### Starting the Application
```
  ██████╗  ██████╗ 
 ██╔════╝ ██╔═══██╗
 ██║  ███╗██║   ██║
 ██║   ██║██║   ██║
 ╚██████╔╝╚██████╔╝
  ╚═════╝  ╚═════╝ 
 
 ██████╗██╗     ██╗
██╔════╝██║     ██║
██║     ██║     ██║
██║     ██║     ██║
╚██████╗███████╗██║
 ╚═════╝╚══════╝╚═╝

╔══════════════════════════╗
║   Application Info      ║
╚══════════════════════════╝
Name: GOCLI ETL Tool
Version: 1.0.0
```

### Progress Bar Example
```
ℹ 📥 Fetching data from api-source
✔ Connected successfully
Downloading from api-source [=========>              ] 45% 2.3 MB/5.1 MB
```

### Menu Example
```
╔══════════════════════════════╗
║       Main Menu             ║
╚══════════════════════════════╝

┌────────┬──────────────────────┐
│ Option │ Description          │
├────────┼──────────────────────┤
│ q      │ 🚪 quit              │
│ 1      │ 💾 state             │
│ 2      │ 🔄 pipelines         │
│ 3      │ 📥 extract           │
│ 4      │ 📤 load              │
│ 5      │ ⚙️ transform          │
│ 6      │ 📋 sequences         │
└────────┴──────────────────────┘

→ Enter option: 
```

### Success Message Example
```
✔ ✅ Pipeline created successfully!
```

### Error Message Example
```
✖ ❌ Failed to load sources: file not found
```

## Benefits

1. **Improved Readability** - Color coding and icons make it easier to scan output
2. **Better Feedback** - Progress bars show operation status in real-time
3. **Professional Appearance** - Modern styling makes the tool more polished
4. **Reduced Cognitive Load** - Visual hierarchy helps users navigate faster
5. **Clear Status** - Icons and colors quickly communicate success/failure

## Configuration

The styling is built-in and requires no configuration. The libraries automatically:
- Detect terminal capabilities
- Adapt to terminal width
- Handle color support gracefully
- Provide fallbacks for limited terminals
