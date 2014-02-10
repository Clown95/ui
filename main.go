// 7 february 2014
package main

import (
	"fmt"
	"os"
	"runtime"
)

func fatalf(format string, args ...interface{}) {
	s := fmt.Sprintf(format, args...)
	_, err := MessageBox(NULL,
		"An internal error has occured:\n" + s,
		os.Args[0],
		MB_OK | MB_ICONERROR)
	if err == nil {
		os.Exit(1)
	}
	panic(fmt.Sprintf("error trying to warn user of internal error: %v\ninternal error:\n%s", err, s))
}

const (
	IDC_BUTTON = 100
	IDC_VARCOMBO = 101
	IDC_FIXCOMBO = 102
)

var varCombo, fixCombo HWND

func wndProc(hwnd HWND, msg uint32, wParam WPARAM, lParam LPARAM) LRESULT {
	switch msg {
	case WM_COMMAND:
		if wParam.LOWORD() == IDC_BUTTON {
			buttonclick := "neither clicked nor double clicked (somehow)"
			if wParam.HIWORD() == BN_CLICKED {
				buttonclick = "clicked"
			} else if wParam.HIWORD() == BN_DOUBLECLICKED {
				buttonclick = "double clicked"
			}

			varText, err := getText(varCombo)
			if err != nil {
				fatalf("error getting variable combo box text: %v", err)
			}

			fixTextWM, err := getText(fixCombo)
			if err != nil {
				fatalf("error getting fixed combo box text with WM_GETTEXT: %v", err)
			}

			fixTextIndex, err := SendMessage(fixCombo, CB_GETCURSEL, 0, 0)
			if err != nil {
				fatalf("error getting fixed combo box current selection: %v", err)
			}
			// TODO get text from index

			MessageBox(hwnd,
				fmt.Sprintf("button state: %s\n" +
					"variable combo box text: %s\n" +
					"fixed combo box text with WM_GETTEXT: %s\n" +
					"fixed combo box current index: %d\n",
					buttonclick, varText, fixTextWM, fixTextIndex),
				"note",
				MB_OK)
		}
		return 0
	case WM_CLOSE:
		err := DestroyWindow(hwnd)
		if err != nil {
			fatalf("error destroying window: %v", err)
		}
		return 0
	case WM_DESTROY:
		err := PostQuitMessage(0)
		if err != nil {
			fatalf("error posting quit message: %v", err)
		}
		return 0
	default:
		return DefWindowProc(hwnd, msg, wParam, lParam)
	}
	fatalf("major bug: forgot a return on wndProc for message %d", msg)
	panic("unreachable")
}

const className = "mainwin"

func main() {
	runtime.LockOSThread()

	hInstance, err := getWinMainhInstance()
	if err != nil {
		fatalf("error getting WinMain hInstance: %v", err)
	}
	nCmdShow, err := getWinMainnCmdShow()
	if err != nil {
		fatalf("error getting WinMain nCmdShow: %v", err)
	}

	icon, err := LoadIcon_ResourceID(NULL, IDI_APPLICATION)
	if err != nil {
		fatalf("error getting window icon: %v", err)
	}
	cursor, err := LoadCursor_ResourceID(NULL, IDC_ARROW)
	if err != nil {
		fatalf("error getting window cursor: %v", err)
	}

	wc := &WNDCLASS{
		LpszClassName:	className,
		LpfnWndProc:		wndProc,
		HInstance:		hInstance,
		HIcon:			icon,
		HCursor:			cursor,
		HbrBackground:	HBRUSH(COLOR_WINDOW + 1),
	}
	_, err = RegisterClass(wc)
	if err != nil {
		fatalf("error registering window class: %v", err)
	}

	hwnd, err := CreateWindowEx(
		WS_EX_OVERLAPPEDWINDOW,
		className, "Main Window",
		WS_OVERLAPPEDWINDOW,
		CW_USEDEFAULT, CW_USEDEFAULT, 320, 240,
		NULL, NULL, hInstance, NULL)
	if err != nil {
		fatalf("error creating window: %v", err)
	}

	_, err = CreateWindowEx(
		0,
		"BUTTON", "Click Me",
		BS_PUSHBUTTON | WS_CHILD | WS_VISIBLE | WS_TABSTOP,
		20, 20, 100, 20,
		hwnd, HMENU(IDC_BUTTON), hInstance, NULL)
	if err != nil {
		fatalf("error creating button: %v", err)
	}

	varCombo, err = CreateWindowEx(
		0,
		"COMBOBOX", "",
		CBS_DROPDOWN | CBS_AUTOHSCROLL | WS_CHILD | WS_VISIBLE | WS_TABSTOP,
		140, 20, 100, 20,
		hwnd, HMENU(IDC_VARCOMBO), hInstance, NULL)
	if err != nil {
		fatalf("error creating variable combo box: %v", err)
	}
	vcItems := []string{"a", "b", "c", "d"}
	for _, v := range vcItems {
		_, err := SendMessage(varCombo, CB_ADDSTRING, 0,
			LPARAMFromString(v))
		if err != nil {
			fatalf("error adding %q to variable combo box: %v", v, err)
		}
	}

	fixCombo, err = CreateWindowEx(
		0,
		"COMBOBOX", "",
		CBS_DROPDOWNLIST | WS_CHILD | WS_VISIBLE | WS_TABSTOP,
		140, 50, 100, 20,
		hwnd, HMENU(IDC_FIXCOMBO), hInstance, NULL)
	if err != nil {
		fatalf("error creating fixed combo box: %v", err)
	}
	fcItems := []string{"e", "f", "g", "h"}
	for _, v := range fcItems {
		_, err := SendMessage(fixCombo, CB_ADDSTRING, 0,
			LPARAMFromString(v))
		if err != nil {
			fatalf("error adding %q to fixed combo box: %v", v, err)
		}
	}

	_, err = ShowWindow(hwnd, nCmdShow)
	if err != nil {
		fatalf("error showing window: %v", err)
	}
	err = UpdateWindow(hwnd)
	if err != nil {
		fatalf("error updating window: %v", err)
	}

	for {
		msg, done, err := GetMessage(NULL, 0, 0)
		if err != nil {
			fatalf("error getting message: %v", err)
		}
		if done {
			break
		}
		_, err = TranslateMessage(msg)
		if err != nil {
			fatalf("error translating message: %v", err)
		}
		_, err = DispatchMessage(msg)
		if err != nil {
			fatalf("error dispatching message: %v", err)
		}
	}
}

