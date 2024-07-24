package main

/*
#include <dlfcn.h>
#include <stdio.h>

void* lib_handle;
void* (*sym_transcoder_alloc)();
int (*sym_transcoder_add_preset)( void* ptr, int width, int height, int bitrate, int fps );
int (*sym_transcoder_init)( void* ptr, int ticks_per_second, int next_segment_index, int target_segment_duration_ms );
int (*sym_transcoder_on_header)( void* ptr, const void* data, size_t data_size );
int (*sym_transcoder_on_video)( void* ptr, long long pts, long long dts, int keyframe, const void* data, size_t data_size );
int (*sym_transcoder_on_eof)( void* ptr );

int (*sym_transcoder_out_queue_size)( void* ptr );
int (*sym_transcoder_out_header)( void* ptr );
int (*sym_transcoder_out_preset)( void* ptr );
long long (*sym_transcoder_out_pts)( void* ptr );
long long (*sym_transcoder_out_dts)( void* ptr );
int (*sym_transcoder_out_keyframe)( void* ptr );
size_t (*sym_transcoder_out_size)( void* ptr );
int (*sym_transcoder_out_segment_index)( void* ptr );
int (*sym_transcoder_out_segment_duration)( void* ptr );
int (*sym_transcoder_out_last_frame_in_segment)( void* ptr );
int (*sym_transcoder_take_frame)( void* ptr, void* out_ptr );


void (*sym_transcoder_free)( void* ptr );


int init_lib(int is_darwin)
{
	const char* libname = is_darwin > 0 ? "./libwdp_transcoder_lite.dylib" : "./libwdp_transcoder_lite.so";
	lib_handle = dlopen(libname, RTLD_NOW);
	if(lib_handle == 0)
	{
		printf("can't load library %s: %s\n", libname, dlerror());
		return -1;
	}

	sym_transcoder_alloc = (void* (*)())dlsym(lib_handle, "transcoder_alloc");
	if(sym_transcoder_alloc == 0)
	{
		printf("can't load symbol transcoder_alloc\n");
		return -1;
	}

	sym_transcoder_add_preset = (int (*)( void* ptr, int width, int height, int bitrate, int fps ))dlsym(lib_handle, "transcoder_add_preset");
	if(sym_transcoder_add_preset == 0)
	{
		printf("can't load symbol transcoder_add_preset\n");
		return -1;
	}

	sym_transcoder_init = (int (*)( void* ptr, int ticks_per_second, int next_segment_index, int target_segment_duration_ms ))dlsym(lib_handle, "transcoder_init");
	if(sym_transcoder_init == 0)
	{
		printf("can't load symbol transcoder_init\n");
		return -1;
	}

	sym_transcoder_on_header = (int (*)( void* ptr, const void* data, size_t data_size ))dlsym(lib_handle, "transcoder_on_header");
	if(sym_transcoder_on_header == 0)
	{
		printf("can't load symbol transcoder_on_header\n");
		return -1;
	}

	sym_transcoder_on_video = (int (*)( void* ptr, long long pts, long long dts, int keyframe, const void* data, size_t data_size ))dlsym(lib_handle, "transcoder_on_video");
	if(sym_transcoder_on_video == 0)
	{
		printf("can't load symbol transcoder_on_video\n");
		return -1;
	}

	sym_transcoder_on_eof = (int (*)( void* ptr ))dlsym(lib_handle, "transcoder_on_eof");
	if(sym_transcoder_on_eof == 0)
	{
		printf("can't load symbol transcoder_on_eof\n");
		return -1;
	}

	sym_transcoder_out_queue_size = (int (*)( void* ptr ))dlsym(lib_handle, "transcoder_out_queue_size");
	if(sym_transcoder_out_queue_size == 0)
	{
		printf("can't load symbol transcoder_out_queue_size\n");
		return -1;
	}

	sym_transcoder_out_preset = (int (*)( void* ptr ))dlsym(lib_handle, "transcoder_out_preset");
	if(sym_transcoder_out_preset == 0)
	{
		printf("can't load symbol transcoder_out_preset\n");
		return -1;
	}

	sym_transcoder_out_header = (int (*)( void* ptr ))dlsym(lib_handle, "transcoder_out_header");
	if(sym_transcoder_out_header == 0)
	{
		printf("can't load symbol transcoder_out_header\n");
		return -1;
	}

	sym_transcoder_out_pts = (long long (*)( void* ptr ))dlsym(lib_handle, "transcoder_out_pts");
	if(sym_transcoder_out_pts == 0)
	{
		printf("can't load symbol transcoder_out_pts\n");
		return -1;
	}

	sym_transcoder_out_dts = (long long (*)( void* ptr ))dlsym(lib_handle, "transcoder_out_dts");
	if(sym_transcoder_out_dts == 0)
	{
		printf("can't load symbol transcoder_out_dts\n");
		return -1;
	}

	sym_transcoder_out_keyframe = (int (*)( void* ptr ))dlsym(lib_handle, "transcoder_out_keyframe");
	if(sym_transcoder_out_keyframe == 0)
	{
		printf("can't load symbol transcoder_out_keyframe\n");
		return -1;
	}

	sym_transcoder_out_size = (size_t (*)( void* ptr ))dlsym(lib_handle, "transcoder_out_size");
	if(sym_transcoder_out_size == 0)
	{
		printf("can't load symbol transcoder_out_size\n");
		return -1;
	}

	sym_transcoder_out_segment_index = (int (*)( void* ptr ))dlsym(lib_handle, "transcoder_out_segment_index");
	if(sym_transcoder_out_segment_index == 0)
	{
		printf("can't load symbol transcoder_out_segment_index\n");
		return -1;
	}

	sym_transcoder_out_segment_duration = (int (*)( void* ptr ))dlsym(lib_handle, "transcoder_out_segment_duration");
	if(sym_transcoder_out_segment_duration == 0)
	{
		printf("can't load symbol transcoder_out_segment_duration\n");
		return -1;
	}

	sym_transcoder_out_last_frame_in_segment = (int (*)( void* ptr ))dlsym(lib_handle, "transcoder_out_last_frame_in_segment");
	if(sym_transcoder_out_last_frame_in_segment == 0)
	{
		printf("can't load symbol transcoder_out_last_frame_in_segment\n");
		return -1;
	}

	sym_transcoder_take_frame = (int (*)( void* ptr, void* out_ptr ))dlsym(lib_handle, "transcoder_take_frame");
	if(sym_transcoder_take_frame == 0)
	{
		printf("can't load symbol transcoder_take_frame\n");
		return -1;
	}

	sym_transcoder_free = (void (*)( void* ptr ))dlsym(lib_handle, "transcoder_free");
	if(sym_transcoder_free == 0)
	{
		printf("can't load symbol transcoder_free\n");
		return -1;
	}

	return 0;
}

void* transcoder_alloc()
{
	if(sym_transcoder_alloc == 0)
	{
		return 0;
	}

	return sym_transcoder_alloc();
}

int transcoder_add_preset( void* ptr, int width, int height, int bitrate, int fps )
{
	if(sym_transcoder_add_preset == 0)
	{
		return -1;
	}

	return sym_transcoder_add_preset(ptr, width, height, bitrate, fps);
}

int transcoder_init( void* ptr, int ticks_per_second, int next_segment_index, int target_segment_duration_ms )
{
	if(sym_transcoder_init == 0)
	{
		return -1;
	}

	return sym_transcoder_init(ptr, ticks_per_second, next_segment_index, target_segment_duration_ms);
}

int transcoder_on_header( void* ptr, const void* data, size_t data_size )
{
	if(sym_transcoder_on_header == 0)
	{
		return -1;
	}

	return sym_transcoder_on_header(ptr, data, data_size);
}

int transcoder_on_video( void* ptr, long long pts, long long dts, int keyframe, const void* data, size_t data_size )
{
	if(sym_transcoder_on_video == 0)
	{
		return -1;
	}

	return sym_transcoder_on_video( ptr, pts, dts, keyframe, data, data_size);
}

int transcoder_on_eof( void* ptr )
{
	if(sym_transcoder_on_eof == 0)
	{
		return -1;
	}

	return sym_transcoder_on_eof(ptr);
}

int transcoder_out_queue_size( void* ptr )
{
	if(sym_transcoder_out_queue_size == 0)
	{
		return -1;
	}

	return sym_transcoder_out_queue_size(ptr);
}

int transcoder_out_preset( void* ptr )
{
	if(sym_transcoder_out_preset == 0)
	{
		return -1;
	}

	return sym_transcoder_out_preset(ptr);
}

int transcoder_out_header( void* ptr )
{
	if(sym_transcoder_out_header == 0)
	{
		return -1;
	}

	return sym_transcoder_out_header(ptr);
}


int transcoder_out_pts( void* ptr )
{
	if(sym_transcoder_out_pts == 0)
	{
		return -1;
	}

	return sym_transcoder_out_pts(ptr);
}

int transcoder_out_dts( void* ptr )
{
	if(sym_transcoder_out_dts == 0)
	{
		return -1;
	}

	return sym_transcoder_out_dts(ptr);
}

int transcoder_out_keyframe( void* ptr )
{
	if(sym_transcoder_out_keyframe == 0)
	{
		return -1;
	}

	return sym_transcoder_out_keyframe(ptr);
}

int transcoder_out_size( void* ptr )
{
	if(sym_transcoder_out_size == 0)
	{
		return -1;
	}

	return sym_transcoder_out_size(ptr);
}

int transcoder_out_segment_index( void* ptr )
{
	if(sym_transcoder_out_segment_index == 0)
	{
		return -1;
	}

	return sym_transcoder_out_segment_index(ptr);
}

int transcoder_out_segment_duration( void* ptr )
{
	if(sym_transcoder_out_segment_duration == 0)
	{
		return -1;
	}

	return sym_transcoder_out_segment_duration(ptr);
}

int transcoder_out_last_frame_in_segment( void* ptr )
{
	if(sym_transcoder_out_last_frame_in_segment == 0)
	{
		return -1;
	}

	return sym_transcoder_out_last_frame_in_segment(ptr);
}


int transcoder_take_frame( void* ptr, void* out_data )
{
	if(sym_transcoder_take_frame == 0)
	{
		return -1;
	}

	return sym_transcoder_take_frame(ptr, out_data);
}

void transcoder_free( void* ptr )
{
	if(sym_transcoder_free == 0)
	{
		return;
	}

	sym_transcoder_free(ptr);
}




*/
import "C"

import (
	"fmt"
	"log"
	"runtime"
	"unsafe"
)

type WdpWrap struct {
}

func newWdpWrap() WdpWrap {
	isDarwin := 0
	if runtime.GOOS == "darwin" {
		isDarwin = 1
	}
	if C.init_lib(C.int(isDarwin)) != 0 {
		log.Fatalf("init_lib failed")
	}

	return WdpWrap{}
}

func (this WdpWrap) Alloc() WdpInstance {
	res := C.transcoder_alloc()
	if res == nil {
		log.Fatalf("can't alloc transcoder")
	}
	return WdpInstance{
		handle: res,
	}
}

type WdpInstance struct {
	handle         unsafe.Pointer
	ticksPerSecond int
}

func (this *WdpInstance) Close() {
	C.transcoder_free(this.handle)
}

func (this *WdpInstance) AddPreset(width int, height int, bitrate int, fps int) error {
	ret := C.transcoder_add_preset(this.handle, C.int(width), C.int(height), C.int(bitrate), C.int(fps))
	if ret != 0 {
		return fmt.Errorf("transcoder_add_preset failed: %d", ret)
	}
	return nil
}

func (this *WdpInstance) Init(ticksPerSecond int, nextSegmentIndex int, targetSegmentDuration int) error {
	ret := C.transcoder_init(this.handle, C.int(ticksPerSecond), C.int(nextSegmentIndex), C.int(targetSegmentDuration))
	if ret != 0 {
		return fmt.Errorf("transcoder_init failed: %d", ret)
	}
	this.ticksPerSecond = ticksPerSecond
	return nil
}

func (this *WdpInstance) OnHeader(data []byte) error {
	ret := C.transcoder_on_header(this.handle, C.CBytes(data), C.size_t(len(data)))
	if ret != 0 {
		return fmt.Errorf("transcoder_on_header failed: %d", ret)
	}
	return nil
}

func (this *WdpInstance) OnVideo(pts int64, dts int64, keyframe bool, data []byte) error {
	keyframeInt := 0
	if keyframe {
		keyframeInt = 1
	}
	ret := C.transcoder_on_video(this.handle, C.longlong(pts), C.longlong(dts), C.int(keyframeInt), C.CBytes(data), C.size_t(len(data)))
	if ret != 0 {
		return fmt.Errorf("transcoder_on_video failed: %d", ret)
	}
	return nil
}

func (this *WdpInstance) OnEof() error {
	ret := C.transcoder_on_eof(this.handle)
	if ret != 0 {
		return fmt.Errorf("transcoder_on_eof failed: %d", ret)
	}
	return nil
}

type WdpPacket struct {
	PresetIndex  int
	SegmentIndex int
	DurationMs   int
	SegmentEnd   bool
	Header       bool
	Keyframe     bool
	Pts          int64
	Dts          int64
	Payload      []byte
}

func (this *WdpInstance) TakePacket() (*WdpPacket, error) {
	queueSize := C.transcoder_out_queue_size(this.handle)
	if queueSize < 0 {
		return nil, fmt.Errorf("transcoder_out_queue_size failed: %d", queueSize)
	}

	if queueSize == 0 {
		return nil, nil
	}

	preset := C.transcoder_out_preset(this.handle)
	if preset < 0 {
		return nil, fmt.Errorf("transcoder_out_preset failed: %d", preset)
	}

	header := C.transcoder_out_header(this.handle)
	if header < 0 {
		return nil, fmt.Errorf("transcoder_out_header failed: %d", header)
	}

	pts := C.transcoder_out_pts(this.handle)
	if pts < 0 {
		return nil, fmt.Errorf("transcoder_out_pts failed: %d", pts)
	}

	dts := C.transcoder_out_dts(this.handle)
	if dts < 0 {
		return nil, fmt.Errorf("transcoder_out_dts failed: %d", dts)
	}

	keyframe := C.transcoder_out_keyframe(this.handle)
	if keyframe < 0 {
		return nil, fmt.Errorf("transcoder_out_keyframe failed: %d", keyframe)
	}

	segmentIndex := C.transcoder_out_segment_index(this.handle)
	if segmentIndex < 0 {
		return nil, fmt.Errorf("transcoder_out_segment_index failed: %d", segmentIndex)
	}

	segmentDuration := C.transcoder_out_segment_duration(this.handle)
	if segmentDuration < 0 {
		return nil, fmt.Errorf("transcoder_out_segment_duration failed: %d", segmentDuration)
	}

	lastFrameInSegment := C.transcoder_out_last_frame_in_segment(this.handle)
	if lastFrameInSegment < 0 {
		return nil, fmt.Errorf("transcoder_out_last_frame_in_segment failed: %d", lastFrameInSegment)
	}

	size := C.transcoder_out_size(this.handle)
	if size < 0 {
		return nil, fmt.Errorf("transcoder_out_size failed: %d", size)
	}

	payload := make([]byte, size)
	payloadC := C.CBytes(payload)

	ret := C.transcoder_take_frame(this.handle, payloadC)
	if ret < 0 {
		return nil, fmt.Errorf("transcoder_take_frame failed: %d", ret)
	}

	payload = C.GoBytes(payloadC, size)

	return &WdpPacket{
		PresetIndex:  int(preset),
		SegmentIndex: int(segmentIndex),
		DurationMs:   int(segmentDuration) * 1000 / this.ticksPerSecond,
		SegmentEnd:   lastFrameInSegment > 0,
		Header:       header > 0,
		Pts:          int64(pts),
		Dts:          int64(dts),
		Keyframe:     keyframe > 0,
		Payload:      payload,
	}, nil

}
