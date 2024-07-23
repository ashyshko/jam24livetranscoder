# Generate ffprobe output

```
ffmpeg -i ~/Movies/bbb-540p30.mp4 -an -t 0:00:30 -c:v libx264 -bf 0 -x264-params repeat-headers=1 ~/Downloads/test-video.mp4

ffprobe -select_streams 0 -print_format compact= -show_frames ~/Downloads/test-video.mp4 >bbb-540p30-new.txt


# ffprobe -select_streams 0 -print_format compact= -show_frames /Users/andrey/Movies/bbb-540p30.mp4 >bbb-540p30.txt
```

# Run locally

```
go build && ./senderstub -i ~/Downloads/test-video.mp4 -c bbb-540p30-new.txt

# go build && ./senderstub -i ~/Movies/bbb-540p30.mp4 -c bbb-540p30.txt
```
