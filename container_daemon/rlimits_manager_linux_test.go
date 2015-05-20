package container_daemon_test

import (
	"github.com/cloudfoundry-incubator/garden"
	"github.com/cloudfoundry-incubator/garden-linux/container_daemon"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("RlimitsManager", func() {
	var rlimits garden.ResourceLimits
	var mgr *container_daemon.RlimitsManager
	// var systemRlimits syscall.Rlimit
	// var prevRlimit *syscall.Rlimit

	// var getSystemRlimits = func() map[int]*syscall.Rlimit {
	// 	rLimitsMap := make(map[int]*syscall.Rlimit)

	// 	for rLimitId := 0; rLimitId <= 14; rLimitId++ {
	// 		rLimitsMap[rLimitId] = new(syscall.Rlimit)
	// 		Expect(syscall.Getrlimit(rLimitId, rLimitsMap[rLimitId])).To(Succeed())
	// 	}

	// 	return rLimitsMap
	// }

	BeforeEach(func() {
		rlimits = garden.ResourceLimits{}
		mgr = new(container_daemon.RlimitsManager)
	})

	Describe("Encode / Decode roundtrip", func() {
		It("preserves the rlimit values", func() {
			var (
				valAs         uint64 = 1
				valCore       uint64 = 2
				valCpu        uint64 = 3
				valData       uint64 = 4
				valFsize      uint64 = 5
				valLocks      uint64 = 6
				valMemlock    uint64 = 7
				valMsgqueue   uint64 = 8
				valNice       uint64 = 9
				valNofile     uint64 = 10
				valNproc      uint64 = 11
				valRss        uint64 = 12
				valRtprio     uint64 = 13
				valSigpending uint64 = 14
				valStack      uint64 = 15
			)

			rlimits := garden.ResourceLimits{
				As:         &valAs,
				Core:       &valCore,
				Cpu:        &valCpu,
				Data:       &valData,
				Fsize:      &valFsize,
				Locks:      &valLocks,
				Memlock:    &valMemlock,
				Msgqueue:   &valMsgqueue,
				Nice:       &valNice,
				Nofile:     &valNofile,
				Nproc:      &valNproc,
				Rss:        &valRss,
				Rtprio:     &valRtprio,
				Sigpending: &valSigpending,
				Stack:      &valStack,
			}

			env := mgr.EncodeEnv(rlimits)
			Expect(env).To(HaveLen(15))

			newRlimits := mgr.DecodeEnv(env)
			Expect(*newRlimits.As).To(Equal(valAs))
			Expect(*newRlimits.Core).To(Equal(valCore))
			Expect(*newRlimits.Cpu).To(Equal(valCpu))
			Expect(*newRlimits.Data).To(Equal(valData))
			Expect(*newRlimits.Fsize).To(Equal(valFsize))
			Expect(*newRlimits.Locks).To(Equal(valLocks))
			Expect(*newRlimits.Memlock).To(Equal(valMemlock))
			Expect(*newRlimits.Msgqueue).To(Equal(valMsgqueue))
			Expect(*newRlimits.Nice).To(Equal(valNice))
			Expect(*newRlimits.Nofile).To(Equal(valNofile))
			Expect(*newRlimits.Nproc).To(Equal(valNproc))
			Expect(*newRlimits.Rss).To(Equal(valRss))
			Expect(*newRlimits.Rtprio).To(Equal(valRtprio))
			Expect(*newRlimits.Sigpending).To(Equal(valSigpending))
			Expect(*newRlimits.Stack).To(Equal(valStack))
		})
	})

	// 	Context("When an error occurs", func() {
	// 		var (
	// 			rLimitValue1    uint64 = 1000
	// 			rLimitStack     uint64 = 100000
	// 			noFileValue     uint64 = 1048999 // this will cause an error for stack rlimit
	// 			previousRlimits map[int]*syscall.Rlimit
	// 		)

	// 		JustBeforeEach(func() {
	// 			rlimits.Cpu = &rLimitValue1
	// 			rlimits.Fsize = &rLimitValue1
	// 			rlimits.Data = &rLimitValue1
	// 			rlimits.Stack = &rLimitStack
	// 			rlimits.Core = &rLimitValue1
	// 			rlimits.Nofile = &noFileValue

	// 			previousRlimits = getSystemRlimits()
	// 		})

	// 		It("rolls back the previous system rlimit values", func() {
	// 			Expect(mgr.Apply(rlimits)).NotTo(Succeed())

	// 			currentRlimits := getSystemRlimits()
	// 			Expect(currentRlimits).To(Equal(previousRlimits))
	// 		})

	// 		It("does not hang in the next apply call", func(done Done) {
	// 			Expect(mgr.Apply(rlimits)).NotTo(Succeed())

	// 			newRlimits := garden.ResourceLimits{
	// 				Cpu:  &rLimitValue1,
	// 				Core: &rLimitValue1,
	// 			}
	// 			Expect(mgr.Apply(newRlimits)).To(Succeed())
	// 			Expect(mgr.Restore()).To(Succeed())

	// 			close(done)
	// 		}, 10)
	// 	})

	// 	Context("when there is no error", func() {
	// 		It("blocks subsequent apply calls, until restore is called", func(done Done) {
	// 			rLimitValue1 := uint64(1000)
	// 			rLimitValue2 := uint64(2000)

	// 			rlimitsA := garden.ResourceLimits{
	// 				Cpu:  &rLimitValue1,
	// 				Core: &rLimitValue1,
	// 			}

	// 			rlimitsB := garden.ResourceLimits{
	// 				Cpu:  &rLimitValue2,
	// 				Core: &rLimitValue2,
	// 			}

	// 			Expect(mgr.Apply(rlimitsA)).To(Succeed())

	// 			applyReturned := make(chan bool)
	// 			restoreReturned := make(chan bool)
	// 			go func(applyReturned, restoreReturned chan bool) {
	// 				defer GinkgoRecover()

	// 				Expect(mgr.Apply(rlimitsB)).To(Succeed())
	// 				close(applyReturned)
	// 				Expect(mgr.Restore()).To(Succeed())
	// 				close(restoreReturned)
	// 			}(applyReturned, restoreReturned)

	// 			Consistently(applyReturned, time.Second).ShouldNot(BeClosed())

	// 			Expect(mgr.Restore()).To(Succeed())
	// 			Eventually(applyReturned, 100*time.Millisecond, 10*time.Millisecond).Should(BeClosed())

	// 			<-restoreReturned

	// 			close(done)
	// 		}, 20)
	// 	})

	// 	Context("CPU limit", func() {
	// 		var rLimitValue uint64 = 1200

	// 		BeforeEach(func() {
	// 			rlimits.Cpu = &rLimitValue
	// 			prevRlimit = new(syscall.Rlimit)

	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_CPU, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_CPU, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous CPU resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_CPU)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate CPU resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_CPU, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous CPU resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_CPU, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).NotTo(Equal(prevRlimit))

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_CPU, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})

	// 	Context("FSIZE limit", func() {
	// 		var rLimitValue uint64 = 1200

	// 		BeforeEach(func() {
	// 			rlimits.Fsize = &rLimitValue
	// 			prevRlimit = new(syscall.Rlimit)

	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_FSIZE, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_FSIZE, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous FSIZE resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_FSIZE)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate FSIZE resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_FSIZE, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous FSIZE resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_FSIZE, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).NotTo(Equal(prevRlimit))

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_FSIZE, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})

	// 	Context("DATA limit", func() {
	// 		var rLimitValue uint64 = 1200

	// 		BeforeEach(func() {
	// 			rlimits.Data = &rLimitValue
	// 			prevRlimit = new(syscall.Rlimit)

	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_DATA, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_DATA, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous DATA resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_DATA)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate DATA resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_DATA, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous DATA resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_DATA, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).NotTo(Equal(prevRlimit))

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_DATA, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})

	// 	Context("STACK limit", func() {
	// 		var rLimitValue uint64

	// 		BeforeEach(func() {
	// 			// Allowed limit is 2^23 (8388608)
	// 			rLimitValue = uint64(math.Pow(2, 23)) + 1000
	// 			rlimits.Stack = &rLimitValue

	// 			prevRlimit = new(syscall.Rlimit)
	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_STACK, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_STACK, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous STACK resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_STACK)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate STACK resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_STACK, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous STACK resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_STACK, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).NotTo(Equal(prevRlimit))

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_STACK, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})

	// 	Context("CORE limit", func() {
	// 		var rLimitValue uint64

	// 		BeforeEach(func() {
	// 			rLimitValue = 9000
	// 			rlimits.Core = &rLimitValue

	// 			prevRlimit = new(syscall.Rlimit)
	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_CORE, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_CORE, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous CORE resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_CORE)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate CORE resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_CORE, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous CORE resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_CORE, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).NotTo(Equal(prevRlimit))

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_CORE, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})

	// 	Context("RSS limit", func() {
	// 		var rLimitValue uint64

	// 		BeforeEach(func() {
	// 			rLimitValue = 9000
	// 			rlimits.Rss = &rLimitValue

	// 			prevRlimit = new(syscall.Rlimit)
	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_RSS, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_RSS, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous RSS resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_RSS)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate RSS resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_RSS, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous RSS resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_RSS, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).NotTo(Equal(prevRlimit))

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_RSS, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})

	// 	Context("NPROC limit", func() {
	// 		var rLimitValue uint64

	// 		BeforeEach(func() {
	// 			rLimitValue = 500
	// 			rlimits.Nproc = &rLimitValue

	// 			prevRlimit = new(syscall.Rlimit)
	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_NPROC, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_NPROC, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous NPROC resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_NPROC)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate NPROC resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_NPROC, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous NPROC resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_NPROC, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).NotTo(Equal(prevRlimit))

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_NPROC, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})

	// 	Context("NOFILE limit", func() {
	// 		var rLimitValue uint64

	// 		BeforeEach(func() {
	// 			rLimitValue = 800
	// 			rlimits.Nofile = &rLimitValue

	// 			prevRlimit = new(syscall.Rlimit)
	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_NOFILE, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_NOFILE, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous NOFILE resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_NOFILE)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate NOFILE resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_NOFILE, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous NOFILE resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_NOFILE, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).NotTo(Equal(prevRlimit))

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_NOFILE, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})

	// 	Context("MEMLOCK limit", func() {
	// 		var rLimitValue uint64

	// 		BeforeEach(func() {
	// 			rLimitValue = 1024
	// 			rlimits.Memlock = &rLimitValue

	// 			prevRlimit = new(syscall.Rlimit)
	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_MEMLOCK, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_MEMLOCK, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous MEMLOCK resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_MEMLOCK)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate MEMLOCK resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_MEMLOCK, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous MEMLOCK resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_MEMLOCK, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).NotTo(Equal(prevRlimit))

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_MEMLOCK, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})

	// 	Context("AS limit", func() {
	// 		var rLimitValue uint64

	// 		BeforeEach(func() {
	// 			rLimitValue = 1024
	// 			rlimits.As = &rLimitValue

	// 			prevRlimit = new(syscall.Rlimit)
	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_AS, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_AS, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous AS resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_AS)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate AS resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_AS, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous AS resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_AS, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).NotTo(Equal(prevRlimit))

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_AS, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})

	// 	Context("LOCKS limit", func() {
	// 		var rLimitValue uint64

	// 		BeforeEach(func() {
	// 			rLimitValue = 500
	// 			rlimits.Locks = &rLimitValue

	// 			prevRlimit = new(syscall.Rlimit)
	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_LOCKS, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_LOCKS, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous LOCKS resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_LOCKS)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate LOCKS resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_LOCKS, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous LOCKS resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_LOCKS, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).NotTo(Equal(prevRlimit))

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_LOCKS, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})

	// 	Context("SIGPENDING limit", func() {
	// 		var rLimitValue uint64

	// 		BeforeEach(func() {
	// 			rLimitValue = 50000
	// 			rlimits.Sigpending = &rLimitValue

	// 			prevRlimit = new(syscall.Rlimit)
	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_SIGPENDING, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_SIGPENDING, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous SIGPENDING resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_SIGPENDING)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate SIGPENDING resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_SIGPENDING, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous SIGPENDING resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_SIGPENDING, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).NotTo(Equal(prevRlimit))

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_SIGPENDING, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})

	// 	Context("MSGQUEUE limit", func() {
	// 		var rLimitValue uint64

	// 		BeforeEach(func() {
	// 			rLimitValue = 500000
	// 			rlimits.Msgqueue = &rLimitValue

	// 			prevRlimit = new(syscall.Rlimit)
	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_MSGQUEUE, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_MSGQUEUE, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous MSGQUEUE resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_MSGQUEUE)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate MSGQUEUE resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_MSGQUEUE, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous MSGQUEUE resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_MSGQUEUE, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).NotTo(Equal(prevRlimit))

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_MSGQUEUE, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})

	// 	Context("NICE limit", func() {
	// 		var rLimitValue uint64

	// 		BeforeEach(func() {
	// 			rLimitValue = 0
	// 			rlimits.Nice = &rLimitValue

	// 			prevRlimit = new(syscall.Rlimit)
	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_NICE, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_NICE, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous NICE resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_NICE)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate NICE resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_NICE, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous NICE resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_NICE, &systemRlimits)).To(Succeed())

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_NICE, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})

	// 	Context("RTPRIO limit", func() {
	// 		var rLimitValue uint64

	// 		BeforeEach(func() {
	// 			rLimitValue = 500
	// 			rlimits.Rtprio = &rLimitValue

	// 			prevRlimit = new(syscall.Rlimit)
	// 			err := syscall.Getrlimit(container_daemon.RLIMIT_RTPRIO, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		AfterEach(func() {
	// 			err := syscall.Setrlimit(container_daemon.RLIMIT_RTPRIO, prevRlimit)
	// 			Expect(err).ToNot(HaveOccurred())
	// 		})

	// 		It("saves the previous RTPRIO resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(mgr.PreviousRLimitValue(container_daemon.RLIMIT_RTPRIO)).To(Equal(prevRlimit))
	// 		})

	// 		It("sets appropriate RTPRIO resource limit", func() {
	// 			Expect(mgr.Apply(rlimits)).To(Succeed())
	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_RTPRIO, &systemRlimits)).To(Succeed())
	// 			Expect(systemRlimits.Cur).To(Equal(rLimitValue))
	// 			Expect(systemRlimits.Max).To(Equal(rLimitValue))
	// 		})

	// 		It("restores the previous RTPRIO resource limit", func() {
	// 			err := mgr.Apply(rlimits)
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_RTPRIO, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).NotTo(Equal(prevRlimit))

	// 			err = mgr.Restore()
	// 			Expect(err).ToNot(HaveOccurred())

	// 			Expect(syscall.Getrlimit(container_daemon.RLIMIT_RTPRIO, &systemRlimits)).To(Succeed())
	// 			Expect(&systemRlimits).To(Equal(prevRlimit))
	// 		})
	// 	})
})
