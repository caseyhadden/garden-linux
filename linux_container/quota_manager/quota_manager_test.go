package quota_manager_test

import (
	"errors"
	"os/exec"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/pivotal-golang/lager/lagertest"

	"github.com/cloudfoundry-incubator/garden"
	"github.com/cloudfoundry-incubator/garden-linux/linux_container/quota_manager"
	"github.com/cloudfoundry/gunk/command_runner/fake_command_runner"
	. "github.com/cloudfoundry/gunk/command_runner/fake_command_runner/matchers"
)

var _ = Describe("btrfs quota manager", func() {
	var fakeRunner *fake_command_runner.FakeCommandRunner
	var logger *lagertest.TestLogger
	var quotaManager *quota_manager.BtrfsQuotaManager
	var containerId string
	var graphRoot string

	BeforeEach(func() {
		fakeRunner = fake_command_runner.New()
		logger = lagertest.NewTestLogger("test")
		graphRoot = "/graph/root/path"
		quotaManager = quota_manager.New(fakeRunner, "/some/mount/point", graphRoot)
		containerId = "some-container"
	})

	Describe("setting quotas", func() {
		limits := garden.DiskLimits{
			ByteSoft: 1,
			ByteHard: 2,

			InodeSoft: 11,
			InodeHard: 12,
		}

		Context("when the subvolume exists", func() {
			BeforeEach(func() {
				fakeRunner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: "btrfs",
						Args: []string{
							"subvolume", "list", graphRoot,
						},
					},
					func(cmd *exec.Cmd) error {
						cmd.Stdout.Write([]byte(
							`ID 11 gen 10 top level 5 path root/path/some/whatever/path
ID 12 gen 10 top level 5 path root/path/some/whatever-1/path/some-container
ID 13 gen 10 top level 5 path root/path/some/whatever-2/path
`,
						))

						return nil
					},
				)
			})

			It("executes qgroup limit with the correct qgroup id", func() {
				err := quotaManager.SetLimits(logger, "some-container", limits)
				Expect(err).ToNot(HaveOccurred())

				Expect(fakeRunner).To(HaveExecutedSerially(
					fake_command_runner.CommandSpec{
						Path: "btrfs",
						Args: []string{
							"qgroup", "limit", "2", "0/12", graphRoot + "/some/whatever-1/path/some-container",
						},
					},
				))
			})

			Context("when blocks are given", func() {
				limits := garden.DiskLimits{
					BlockSoft: 10,
					BlockHard: 20,
				}

				It("executes qgroup limit with them converted to bytes", func() {
					err := quotaManager.SetLimits(logger, containerId, limits)

					Expect(err).ToNot(HaveOccurred())

					Expect(fakeRunner).To(HaveExecutedSerially(
						fake_command_runner.CommandSpec{
							Path: "btrfs",
							Args: []string{
								"qgroup", "limit", "20480", "0/12", graphRoot + "/some/whatever-1/path/some-container"},
						},
					))
				})
			})

			Context("when executing qgroup limit fails", func() {
				nastyError := errors.New("oh no!")

				BeforeEach(func() {
					fakeRunner.WhenRunning(
						fake_command_runner.CommandSpec{
							Path: "btrfs",
						}, func(*exec.Cmd) error {
							return nastyError
						},
					)
				})

				It("returns the error", func() {
					err := quotaManager.SetLimits(logger, containerId, limits)
					Expect(err).To(MatchError("quota_manager: failed to apply limit: oh no!"))
				})
			})
		})

		Context("when the desired subvolume cannot be found", func() {
			var btrfsSubvolResponse []byte

			JustBeforeEach(func() {
				fakeRunner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: "btrfs",
						Args: []string{
							"subvolume", "list", "/some/mount/point",
						},
					},
					func(cmd *exec.Cmd) error {
						cmd.Stdout.Write(btrfsSubvolResponse)
						return nil
					},
				)
			})

			Context("When there are no subvolumes", func() {
				BeforeEach(func() {
					btrfsSubvolResponse = []byte("")
				})

				It("returns an error", func() {
					err := quotaManager.SetLimits(logger, containerId, limits)
					Expect(err).To(MatchError(ContainSubstring("quota_manager: subvolume not found")))
				})
			})

			Context("when there are subvolumes and the container subvolume does not exist", func() {

				BeforeEach(func() {
					btrfsSubvolResponse = []byte(
						`ID 11 gen 10 top level 5 path some/whatever/path
ID 13 gen 10 top level 5 path some/whatever-2/path
`,
					)
				})

				It("returns an error", func() {
					err := quotaManager.SetLimits(logger, containerId, limits)
					Expect(err).To(MatchError(ContainSubstring("quota_manager: subvolume not found")))
				})
			})
		})

		Context("when quotas are disabled", func() {
			BeforeEach(func() {
				quotaManager.Disable()
			})

			It("runs nothing", func() {
				err := quotaManager.SetLimits(logger, containerId, limits)

				Expect(err).ToNot(HaveOccurred())

				for _, cmd := range fakeRunner.ExecutedCommands() {
					Expect(cmd.Path).ToNot(Equal("btrfs"))
				}
			})
		})
	})

	Describe("getting quotas limits", func() {
		BeforeEach(func() {
			fakeRunner.WhenRunning(
				fake_command_runner.CommandSpec{
					Path: "btrfs",
					Args: []string{
						"subvolume", "list", graphRoot,
					},
				},
				func(cmd *exec.Cmd) error {
					cmd.Stdout.Write([]byte(
						`ID 11 gen 10 top level 5 path path/whatever/volume1
ID 12 gen 10 top level 5 path path/whatever/volume2
ID 13 gen 10 top level 5 path path/whatever/volume3
ID 14 gen 10 top level 5 path path/whatever/some-container
`,
					))

					return nil
				},
			)
		})

		It("gets current quotas using btrfs", func() {
			fakeRunner.WhenRunning(
				fake_command_runner.CommandSpec{}, func(cmd *exec.Cmd) error {
					cmd.Stdout.Write([]byte("1000000\n"))
					return nil
				},
			)

			limits, err := quotaManager.GetLimits(logger, containerId)
			Expect(err).ToNot(HaveOccurred())

			Expect(limits.ByteSoft).To(Equal(uint64(1000000)))
			Expect(limits.ByteHard).To(Equal(uint64(1000000)))
		})

		Context("when getting quota using btrfs fails", func() {
			disaster := errors.New("oh no!")

			BeforeEach(func() {
				fakeRunner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: "sh",
					}, func(cmd *exec.Cmd) error {
						return disaster
					},
				)
			})

			It("returns the error", func() {
				_, err := quotaManager.GetLimits(logger, containerId)
				Expect(err).To(MatchError(ContainSubstring("quota_manager: failed to get limit")))
			})
		})

		Context("when getting quota using btrfs spews out malformed results", func() {
			BeforeEach(func() {
				fakeRunner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: "sh",
					}, func(cmd *exec.Cmd) error {
						cmd.Stdout.Write([]byte("Oops\n"))
						return nil
					},
				)
			})

			It("returns the error", func() {
				_, err := quotaManager.GetLimits(logger, containerId)
				Expect(err).To(MatchError(ContainSubstring("quota_manager: failed to parse result")))
			})
		})

		Context("when the output of repquota is malformed", func() {
			It("returns an error", func() {
				fakeRunner.WhenRunning(
					fake_command_runner.CommandSpec{
						Path: "btrfs",
						// Args: []string{"/some/mount/point", "1234"},
					}, func(cmd *exec.Cmd) error {
						cmd.Stdout.Write([]byte("abc\n"))

						return nil
					},
				)

				_, err := quotaManager.GetLimits(logger, containerId)
				Expect(err).To(HaveOccurred())
			})
		})

		Context("when quotas are disabled", func() {
			BeforeEach(func() {
				quotaManager.Disable()
			})

			It("runs nothing", func() {
				limits, err := quotaManager.GetLimits(logger, containerId)
				Expect(err).ToNot(HaveOccurred())

				Expect(limits).To(BeZero())

				for _, cmd := range fakeRunner.ExecutedCommands() {
					Expect(cmd.Path).ToNot(Equal("btrfs"))
				}
			})
		})
	})

	//PDescribe("getting usage", func() {
	//	It("executes repquota in the root path", func() {
	//		fakeRunner.WhenRunning(
	//			fake_command_runner.CommandSpec{
	//				Path: "/root/path/repquota",
	//				Args: []string{"/some/mount/point", "1234"},
	//			}, func(cmd *exec.Cmd) error {
	//				cmd.Stdout.Write([]byte("1234 111 222 333 444 555 666 777 888\n"))

	//				return nil
	//			},
	//		)

	//		limits, err := quotaManager.GetUsage(logger, 1234)
	//		Expect(err).ToNot(HaveOccurred())

	//		Expect(limits.BytesUsed).To(Equal(uint64(111)))
	//		Expect(limits.InodesUsed).To(Equal(uint64(555)))
	//	})

	//	Context("when repquota fails", func() {
	//		disaster := errors.New("oh no!")

	//		BeforeEach(func() {
	//			fakeRunner.WhenRunning(
	//				fake_command_runner.CommandSpec{
	//					Path: "/root/path/repquota",
	//					Args: []string{"/some/mount/point", "1234"},
	//				}, func(cmd *exec.Cmd) error {
	//					return disaster
	//				},
	//			)
	//		})

	//		It("returns the error", func() {
	//			_, err := quotaManager.GetUsage(logger, 1234)
	//			Expect(err).To(Equal(disaster))
	//		})
	//	})

	//	Context("when the output of repquota is malformed", func() {
	//		It("returns an error", func() {
	//			fakeRunner.WhenRunning(
	//				fake_command_runner.CommandSpec{
	//					Path: "/root/path/repquota",
	//					Args: []string{"/some/mount/point", "1234"},
	//				}, func(cmd *exec.Cmd) error {
	//					cmd.Stdout.Write([]byte("abc\n"))

	//					return nil
	//				},
	//			)

	//			_, err := quotaManager.GetUsage(logger, 1234)
	//			Expect(err).To(HaveOccurred())
	//		})
	//	})

	//	Context("when quotas are disabled", func() {
	//		BeforeEach(func() {
	//			quotaManager.Disable()
	//		})

	//		It("runs nothing", func() {
	//			usage, err := quotaManager.GetUsage(logger, 1234)
	//			Expect(err).ToNot(HaveOccurred())

	//			Expect(usage).To(BeZero())

	//			for _, cmd := range fakeRunner.ExecutedCommands() {
	//				Expect(cmd.Path).ToNot(Equal("btrfs"))
	//			}
	//		})
	//	})
	//})

	PDescribe("getting the mount point", func() {
		It("returns the mount point of the container depot", func() {
			Expect(quotaManager.MountPoint()).To(Equal("/some/mount/point"))
		})
	})
})
