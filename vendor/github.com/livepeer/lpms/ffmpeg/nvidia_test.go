// +build nvidia

package ffmpeg

import (
	"fmt"
	"os"
	"testing"
)

func TestNvidia_Transcoding(t *testing.T) {
	// Various Nvidia GPU tests for encoding + decoding
	// XXX what is missing is a way to verify these are *actually* running on GPU!

	_, dir := setupTest(t)
	defer os.RemoveAll(dir)

	var err error
	fname := "../transcoder/test.ts"
	oname := dir + "/out.ts"
	prof := P240p30fps16x9

	// hw enc, sw dec
	err = Transcode2(&TranscodeOptionsIn{
		Fname: fname,
		Accel: Nvidia,
	}, []TranscodeOptions{
		TranscodeOptions{
			Oname:   oname,
			Profile: prof,
			Accel:   Software,
		},
	})
	if err != nil {
		t.Error(err)
	}

	// sw dec, hw enc
	err = Transcode2(&TranscodeOptionsIn{
		Fname: fname,
		Accel: Software,
	}, []TranscodeOptions{
		TranscodeOptions{
			Oname:   oname,
			Profile: prof,
			Accel:   Nvidia,
		},
	})
	if err != nil {
		t.Error(err)
	}

	// hw enc + dec
	err = Transcode2(&TranscodeOptionsIn{
		Fname: fname,
		Accel: Nvidia,
	}, []TranscodeOptions{
		TranscodeOptions{
			Oname:   oname,
			Profile: prof,
			Accel:   Nvidia,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestNvidia_Pixfmts(t *testing.T) {

	run, dir := setupTest(t)
	defer os.RemoveAll(dir)

	oname := dir + "/out.ts"
	prof := P240p30fps16x9

	// check valid and invalid pixel formats
	cmd := `
    set -eux
    cd "$0"
    cp "$1/../transcoder/test.ts" test.ts

    # sanity check original input type is 420p
    ffprobe -loglevel warning test.ts  -show_streams -select_streams v | grep pix_fmt=yuv420p

    # generate invalid 422p type
    ffmpeg -loglevel warning -i test.ts -an -c:v libx264 -pix_fmt yuv422p -vframes 1 out422p.mp4
    ffprobe -loglevel warning out422p.mp4  -show_streams -select_streams v | grep pix_fmt=yuv422p

    # generate valid 444p type
    ffmpeg -loglevel warning -i test.ts -an -c:v libx264 -pix_fmt yuv444p -vframes 1 out444p.mp4
    ffprobe -loglevel warning out444p.mp4  -show_streams -select_streams v | grep pix_fmt=yuv444p
  `
	run(cmd)

	// check invalid pixel format
	err := Transcode2(&TranscodeOptionsIn{
		Fname: dir + "/out422p.mp4",
		Accel: Nvidia,
	}, []TranscodeOptions{
		TranscodeOptions{
			Oname:   oname,
			Profile: prof,
			Accel:   Nvidia,
		},
	})
	if err == nil || err.Error() != "Unsupported input pixel format" {
		t.Error(err)
	}

	// Software decode an invalid GPU pixfmt then use hw encoding
	// XXX need to convert pixfmt in software first before uploading
	/*
		err = Transcode2(&TranscodeOptionsIn{
			Fname: dir + "/out422p.mp4",
			Accel: Software,
		}, []TranscodeOptions{
			TranscodeOptions{
				Oname:   oname,
				Profile: prof,
				Accel:   Nvidia,
			},
		})
		if err != nil {
			t.Error(err)
		}
	*/

	// check different type of valid pixfmt.
	// Some cards (eg, Tesla K80 on GCP) does not have support for these!
	err = Transcode2(&TranscodeOptionsIn{
		Fname: dir + "/out444p.mp4",
		Accel: Nvidia,
	}, []TranscodeOptions{
		TranscodeOptions{
			Oname:   oname,
			Profile: prof,
			Accel:   Nvidia,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestNvidia_Transcoding_Multiple(t *testing.T) {

	// Tests multiple encoding profiles.
	// May be skipped in 'short' mode.

	if testing.Short() {
		t.Skip("Skipping encoding multiple profiles")
	}

	_, dir := setupTest(t)
	defer os.RemoveAll(dir)

	fname := "../transcoder/test.ts"
	prof := P240p30fps16x9

	// hw enc + dec, multiple
	mkoname := func(i int) string { return fmt.Sprintf("%s/%d.ts", dir, i) }
	err := Transcode2(&TranscodeOptionsIn{
		Fname: fname,
		Accel: Software,
	}, []TranscodeOptions{
		TranscodeOptions{
			Oname:   mkoname(0),
			Profile: prof,
			Accel:   Nvidia,
		},
		TranscodeOptions{
			Oname:   mkoname(1),
			Profile: prof,
			Accel:   Nvidia,
		},
		TranscodeOptions{
			Oname:   mkoname(2),
			Profile: prof,
			Accel:   Nvidia,
		},
		TranscodeOptions{
			Oname:   mkoname(3),
			Profile: prof,
			Accel:   Nvidia,
		},
		TranscodeOptions{
			Oname:   mkoname(4),
			Profile: prof,
			Accel:   Nvidia,
		},
	})
	if err != nil {
		t.Error(err)
	}
}

func TestNvidia_Devices(t *testing.T) {

	// XXX need to verify these are running on the correct GPU
	//     not just that the code runs

	device := os.Getenv("GPU_DEVICE")
	if device == "" {
		t.Skip("Skipping device specific tests; no GPU_DEVICE set")
	}

	_, dir := setupTest(t)
	defer os.RemoveAll(dir)

	var err error
	fname := "../transcoder/test.ts"
	oname := dir + "/out.ts"
	prof := P240p30fps16x9

	// hw enc, sw dec
	err = Transcode2(&TranscodeOptionsIn{
		Fname:  fname,
		Accel:  Nvidia,
		Device: device,
	}, []TranscodeOptions{
		TranscodeOptions{
			Oname:   oname,
			Profile: prof,
			Accel:   Software,
		},
	})
	if err != nil {
		t.Error(err)
	}

	// sw dec, hw enc
	err = Transcode2(&TranscodeOptionsIn{
		Fname: fname,
		Accel: Software,
	}, []TranscodeOptions{
		TranscodeOptions{
			Oname:   oname,
			Profile: prof,
			Accel:   Nvidia,
			Device:  device,
		},
	})
	if err != nil {
		t.Error(err)
	}

	// hw enc + dec
	err = Transcode2(&TranscodeOptionsIn{
		Fname:  fname,
		Accel:  Nvidia,
		Device: device,
	}, []TranscodeOptions{
		TranscodeOptions{
			Oname:   oname,
			Profile: prof,
			Accel:   Nvidia,
		},
	})
	if err != nil {
		t.Error(err)
	}

	// hw enc + hw dec, separate devices
	err = Transcode2(&TranscodeOptionsIn{
		Fname:  fname,
		Accel:  Nvidia,
		Device: "0",
	}, []TranscodeOptions{
		TranscodeOptions{
			Oname:   oname,
			Profile: prof,
			Accel:   Nvidia,
			Device:  "1",
		},
	})
	if err != ErrTranscoderInp {
		t.Error(err)
	}

	// invalid device for decoding
	err = Transcode2(&TranscodeOptionsIn{
		Fname:  fname,
		Accel:  Nvidia,
		Device: "9999",
	}, []TranscodeOptions{
		TranscodeOptions{
			Oname:   oname,
			Profile: prof,
			Accel:   Software,
		},
	})
	if err == nil || err.Error() != "Unknown error occurred" {
		t.Error(fmt.Errorf(fmt.Sprintf("\nError being: '%v'\n", err)))
	}

	// invalid device for encoding
	err = Transcode2(&TranscodeOptionsIn{
		Fname: fname,
		Accel: Software,
	}, []TranscodeOptions{
		TranscodeOptions{
			Oname:   oname,
			Profile: prof,
			Accel:   Nvidia,
			Device:  "9999",
		},
	})
	if err == nil || err.Error() != "Unknown error occurred" {
		t.Error(fmt.Errorf(fmt.Sprintf("\nError being: '%v'\n", err)))
	}
}
