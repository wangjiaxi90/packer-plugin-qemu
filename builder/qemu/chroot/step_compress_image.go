package chroot

import (
	"context"
	"fmt"
	"runtime"

	"github.com/hashicorp/packer-plugin-sdk/multistep"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
)

// StepCompressImage create and compress the final image
type StepCompressImage struct {
}

func (s *StepCompressImage) Run(_ context.Context, state multistep.StateBag) multistep.StepAction {
	ui := state.Get("ui").(packersdk.Ui)
	imagePath := state.Get("image_path").(string)
	rawImage := state.Get("rawImage").(string)

	ui.Say("Compressing image...")

	if runtime.GOARCH == "arm64" {
		// https://bugzilla.redhat.com/show_bug.cgi?id=1969848
		if _, err := RunCommand(state, fmt.Sprintf("qemu-img convert -m 1 -f raw -O qcow2 -c %s %s", rawImage, imagePath)); err != nil {
			return Halt(state, err)
		}
	} else {
		if _, err := RunCommand(state, fmt.Sprintf("qemu-img convert -f raw -O qcow2 -c %s %s", rawImage, imagePath)); err != nil {
			return Halt(state, err)
		}
	}

	return multistep.ActionContinue
}

func (s *StepCompressImage) Cleanup(state multistep.StateBag) {}