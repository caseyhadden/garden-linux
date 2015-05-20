package lifecycle_test

import (
	"io"

	"github.com/cloudfoundry-incubator/garden"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gbytes"
)

var _ = Describe("Resource limits", func() {
	var (
		container           garden.Container
		privilegedContainer bool
	)

	JustBeforeEach(func() {
		var err error

		client = startGarden()

		container, err = client.Create(garden.ContainerSpec{
			Privileged: privilegedContainer,
		})
		Expect(err).ToNot(HaveOccurred())
	})

	AfterEach(func() {
		err := client.Destroy(container.Handle())
		Expect(err).ToNot(HaveOccurred())
	})

	Context("when setting all rlimits to minimum values", func() {
		It("succeeds", func(done Done) {
			var (
				val0 uint64 = 0
				// Number of open files
				valNofile uint64 = 4
				// Memory limits
				valAs    uint64 = 4194304
				valData  uint64 = 8192
				valStack uint64 = 11264
			)

			rlimits := garden.ResourceLimits{
				// Memory limits
				As:    &valAs,    // size of processe's virtual memory
				Data:  &valData,  // data segment size
				Stack: &valStack, // stack segemtn size
				// Number of open files
				Nofile: &valNofile, // number of open file descriptors
				// Can be zero
				Core:       &val0, // core file size - 0 for disabled
				Cpu:        &val0, // cpu time in seconds
				Fsize:      &val0, // total size of created files
				Locks:      &val0, // number of file system locks
				Memlock:    &val0, // number of memory locks
				Msgqueue:   &val0, // size of message queue (in bytes)
				Nice:       &val0, // nice max priority (applyied as 20 - limit)
				Nproc:      &val0, // number of spawned processes
				Rss:        &val0, // number of pages resident in RAM
				Rtprio:     &val0, // limit of real-time priority
				Sigpending: &val0, // number of signals that are queued for the process
			}

			proc, err := container.Run(
				garden.ProcessSpec{
					Path:   "echo",
					Args:   []string{"Hello world"},
					User:   "root",
					Limits: rlimits,
				},
				garden.ProcessIO{
					Stdout: GinkgoWriter,
					Stderr: GinkgoWriter,
				},
			)
			Expect(err).ToNot(HaveOccurred())

			Expect(proc.Wait()).To(Equal(0))

			close(done)
		}, 10)
	})

	Describe("Specific resource limits", func() {
		Context("CPU rlimit", func() {
			Context("with a privileged container", func() {
				BeforeEach(func() {
					privilegedContainer = true
				})

				It("rlimits can be set", func() {
					var cpu uint64 = 9000
					stdout := gbytes.NewBuffer()

					process, err := container.Run(garden.ProcessSpec{
						Path: "sh",
						User: "root",
						Args: []string{"-c", "ulimit -t"},
						Limits: garden.ResourceLimits{
							Cpu: &cpu,
						},
					}, garden.ProcessIO{Stdout: io.MultiWriter(stdout, GinkgoWriter), Stderr: GinkgoWriter})
					Expect(err).ToNot(HaveOccurred())

					Eventually(stdout).Should(gbytes.Say("9000"))
					Expect(process.Wait()).To(Equal(0))
				})
			})

			Context("with a non-privileged container", func() {
				BeforeEach(func() {
					privilegedContainer = false
				})

				It("rlimits can be set", func() {
					var cpu uint64 = 9000
					stdout := gbytes.NewBuffer()

					process, err := container.Run(garden.ProcessSpec{
						Path: "sh",
						User: "root",
						Args: []string{"-c", "ulimit -t"},
						Limits: garden.ResourceLimits{
							Cpu: &cpu,
						},
					}, garden.ProcessIO{Stdout: io.MultiWriter(stdout, GinkgoWriter), Stderr: GinkgoWriter})
					Expect(err).ToNot(HaveOccurred())

					Eventually(stdout).Should(gbytes.Say("9000"))
					Expect(process.Wait()).To(Equal(0))
				})
			})
		})

		Context("FSIZE rlimit", func() {
			Context("with a privileged container", func() {
				BeforeEach(func() {
					privilegedContainer = true
				})

				It("rlimits can be set", func() {
					var fsize uint64 = 4194304
					stdout := gbytes.NewBuffer()

					process, err := container.Run(garden.ProcessSpec{
						Path: "sh",
						User: "root",
						Args: []string{"-c", "ulimit -f"},
						Limits: garden.ResourceLimits{
							Fsize: &fsize,
						},
					}, garden.ProcessIO{Stdout: io.MultiWriter(stdout, GinkgoWriter), Stderr: GinkgoWriter})
					Expect(err).ToNot(HaveOccurred())

					Eventually(stdout).Should(gbytes.Say("8192"))
					Expect(process.Wait()).To(Equal(0))
				})
			})

			Context("with a non-privileged container", func() {
				BeforeEach(func() {
					privilegedContainer = false
				})

				It("rlimits can be set", func() {
					var fsize uint64 = 4194304
					stdout := gbytes.NewBuffer()

					process, err := container.Run(garden.ProcessSpec{
						Path: "sh",
						User: "root",
						Args: []string{"-c", "ulimit -f"},
						Limits: garden.ResourceLimits{
							Fsize: &fsize,
						},
					}, garden.ProcessIO{Stdout: io.MultiWriter(stdout, GinkgoWriter), Stderr: GinkgoWriter})
					Expect(err).ToNot(HaveOccurred())

					Eventually(stdout).Should(gbytes.Say("8192"))
					Expect(process.Wait()).To(Equal(0))
				})
			})
		})

		Context("NOFILE rlimit", func() {
			Context("with a privileged container", func() {
				BeforeEach(func() {
					privilegedContainer = true
				})

				It("rlimits can be set", func() {
					var nofile uint64 = 524288
					stdout := gbytes.NewBuffer()

					process, err := container.Run(garden.ProcessSpec{
						Path: "sh",
						User: "root",
						Args: []string{"-c", "ulimit -n"},
						Limits: garden.ResourceLimits{
							Nofile: &nofile,
						},
					}, garden.ProcessIO{Stdout: io.MultiWriter(stdout, GinkgoWriter), Stderr: GinkgoWriter})
					Expect(err).ToNot(HaveOccurred())

					Eventually(stdout).Should(gbytes.Say("524288"))
					Expect(process.Wait()).To(Equal(0))
				})
			})

			Context("with a non-privileged container", func() {
				BeforeEach(func() {
					privilegedContainer = false
				})

				It("rlimits can be set", func() {
					var nofile uint64 = 524288
					stdout := gbytes.NewBuffer()
					process, err := container.Run(garden.ProcessSpec{
						Path: "sh",
						User: "root",
						Args: []string{"-c", "ulimit -n"},
						Limits: garden.ResourceLimits{
							Nofile: &nofile,
						},
					}, garden.ProcessIO{Stdout: io.MultiWriter(stdout, GinkgoWriter), Stderr: GinkgoWriter})
					Expect(err).ToNot(HaveOccurred())

					Eventually(stdout).Should(gbytes.Say("524288"))
					Expect(process.Wait()).To(Equal(0))
				})
			})
		})
	})
})
