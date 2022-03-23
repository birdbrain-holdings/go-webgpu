package wgpu

/*

#include <stdlib.h>

#include "./lib/wgpu.h"

extern void requestAdapterCallback_cgo(WGPURequestAdapterStatus status,
                                WGPUAdapter adapter, char const *message,
                                void *userdata);

*/
import "C"

import (
	"errors"
	"runtime/cgo"
	"unsafe"
)

type AdapterExtras struct {
	BackendType BackendType
}

type RequestAdapterOptions struct {
	CompatibleSurface    *Surface
	PowerPreference      PowerPreference
	ForceFallbackAdapter bool

	// ChainedStruct -> WGPUAdapterExtras
	AdapterExtras *AdapterExtras
}

func RequestAdapter(options *RequestAdapterOptions) (*Adapter, error) {
	var opts C.WGPURequestAdapterOptions

	if options != nil {
		if options.CompatibleSurface != nil {
			opts.compatibleSurface = options.CompatibleSurface.ref
		}
		opts.powerPreference = C.WGPUPowerPreference(options.PowerPreference)
		opts.forceFallbackAdapter = C.bool(options.ForceFallbackAdapter)

		if options.AdapterExtras != nil {
			adapterExtras := (*C.WGPUAdapterExtras)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUAdapterExtras{}))))
			defer C.free(unsafe.Pointer(adapterExtras))

			adapterExtras.chain.next = nil
			adapterExtras.chain.sType = C.WGPUSType_AdapterExtras
			adapterExtras.backend = C.WGPUBackendType(options.AdapterExtras.BackendType)

			opts.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(adapterExtras))
		}
	}

	var status RequestAdapterStatus
	var adapter *Adapter

	var cb requestAdapterCB = func(s RequestAdapterStatus, a *Adapter, _ string) {
		status = s
		adapter = a
	}
	handle := cgo.NewHandle(cb)
	C.wgpuInstanceRequestAdapter(nil, &opts, C.WGPURequestAdapterCallback(C.requestAdapterCallback_cgo), unsafe.Pointer(&handle))

	if status != RequestAdapterStatus_Success {
		return nil, errors.New("failed to request adapter")
	}
	return adapter, nil
}

type SurfaceDescriptorFromWindowsHWND struct {
	Hinstance unsafe.Pointer
	Hwnd      unsafe.Pointer
}

type SurfaceDescriptorFromXcbWindow struct {
	Connection unsafe.Pointer
	Window     uint32
}

type SurfaceDescriptorFromXlibWindow struct {
	Display unsafe.Pointer
	Window  uint32
}

type SurfaceDescriptorFromMetalLayer struct {
	Layer unsafe.Pointer
}

type SurfaceDescriptorFromWaylandSurface struct {
	Display unsafe.Pointer
	Surface unsafe.Pointer
}

type SurfaceDescriptorFromAndroidNativeWindow struct {
	Window unsafe.Pointer
}

type SurfaceDescriptor struct {
	Label string

	// ChainedStruct -> WGPUSurfaceDescriptorFromWindowsHWND
	WindowsHWND *SurfaceDescriptorFromWindowsHWND

	// ChainedStruct -> WGPUSurfaceDescriptorFromXcbWindow
	XcbWindow *SurfaceDescriptorFromXcbWindow

	// ChainedStruct -> WGPUSurfaceDescriptorFromXlibWindow
	XlibWindow *SurfaceDescriptorFromXlibWindow

	// ChainedStruct -> WGPUSurfaceDescriptorFromMetalLayer
	MetalLayer *SurfaceDescriptorFromMetalLayer

	// ChainedStruct -> WGPUSurfaceDescriptorFromWaylandSurface
	WaylandSurface *SurfaceDescriptorFromWaylandSurface

	// ChainedStruct -> WGPUSurfaceDescriptorFromAndroidNativeWindow
	AndroidNativeWindow *SurfaceDescriptorFromAndroidNativeWindow
}

func CreateSurface(descriptor *SurfaceDescriptor) *Surface {
	var desc C.WGPUSurfaceDescriptor

	if descriptor != nil {
		if descriptor.Label != "" {
			label := C.CString(descriptor.Label)
			defer C.free(unsafe.Pointer(label))

			desc.label = label
		}

		if descriptor.WindowsHWND != nil {
			windowsHWND := (*C.WGPUSurfaceDescriptorFromWindowsHWND)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceDescriptorFromWindowsHWND{}))))
			defer C.free(unsafe.Pointer(windowsHWND))

			windowsHWND.chain.next = nil
			windowsHWND.chain.sType = C.WGPUSType_SurfaceDescriptorFromWindowsHWND
			windowsHWND.hinstance = descriptor.WindowsHWND.Hinstance
			windowsHWND.hwnd = descriptor.WindowsHWND.Hwnd

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(windowsHWND))
		}

		if descriptor.XcbWindow != nil {
			xcbWindow := (*C.WGPUSurfaceDescriptorFromXcbWindow)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceDescriptorFromXcbWindow{}))))
			defer C.free(unsafe.Pointer(xcbWindow))

			xcbWindow.chain.next = nil
			xcbWindow.chain.sType = C.WGPUSType_SurfaceDescriptorFromXcbWindow
			xcbWindow.connection = descriptor.XcbWindow.Connection
			xcbWindow.window = C.uint32_t(descriptor.XcbWindow.Window)

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(xcbWindow))
		}

		if descriptor.XlibWindow != nil {
			xlibWindow := (*C.WGPUSurfaceDescriptorFromXlibWindow)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceDescriptorFromXlibWindow{}))))
			defer C.free(unsafe.Pointer(xlibWindow))

			xlibWindow.chain.next = nil
			xlibWindow.chain.sType = C.WGPUSType_SurfaceDescriptorFromXlibWindow
			xlibWindow.display = descriptor.XlibWindow.Display
			xlibWindow.window = C.uint32_t(descriptor.XlibWindow.Window)

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(xlibWindow))
		}

		if descriptor.MetalLayer != nil {
			metalLayer := (*C.WGPUSurfaceDescriptorFromMetalLayer)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceDescriptorFromMetalLayer{}))))
			defer C.free(unsafe.Pointer(metalLayer))

			metalLayer.chain.next = nil
			metalLayer.chain.sType = C.WGPUSType_SurfaceDescriptorFromMetalLayer
			metalLayer.layer = descriptor.MetalLayer.Layer

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(metalLayer))
		}

		if descriptor.WaylandSurface != nil {
			waylandSurface := (*C.WGPUSurfaceDescriptorFromWaylandSurface)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceDescriptorFromWaylandSurface{}))))
			defer C.free(unsafe.Pointer(waylandSurface))

			waylandSurface.chain.next = nil
			waylandSurface.chain.sType = C.WGPUSType_SurfaceDescriptorFromWaylandSurface
			waylandSurface.display = descriptor.WaylandSurface.Display
			waylandSurface.surface = descriptor.WaylandSurface.Surface

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(waylandSurface))
		}

		if descriptor.AndroidNativeWindow != nil {
			androidNativeWindow := (*C.WGPUSurfaceDescriptorFromAndroidNativeWindow)(C.malloc(C.size_t(unsafe.Sizeof(C.WGPUSurfaceDescriptorFromAndroidNativeWindow{}))))
			defer C.free(unsafe.Pointer(androidNativeWindow))

			androidNativeWindow.chain.next = nil
			androidNativeWindow.chain.sType = C.WGPUSType_SurfaceDescriptorFromAndroidNativeWindow
			androidNativeWindow.window = descriptor.AndroidNativeWindow.Window

			desc.nextInChain = (*C.WGPUChainedStruct)(unsafe.Pointer(androidNativeWindow))
		}
	}

	ref := C.wgpuInstanceCreateSurface(nil, &desc)
	if ref == nil {
		panic("Failed to acquire Surface")
	}
	return &Surface{ref}
}
