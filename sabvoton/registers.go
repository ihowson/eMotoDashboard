package sabvoton

type RegisterType int

const (
	Int16 RegisterType = iota
	UInt16
	Float16
	DFloat16
)

type FloatPrecision float32

const (
	Single FloatPrecision = 54.6
	Double FloatPrecision = 6553.0
)

type RegisterUInt16 struct {
	Address uint16
	Sample  int
}

type RegisterFloat16 struct {
	Address   uint16
	Type      RegisterType
	Precision FloatPrecision
	Sample    float32
}

var RegisterDCCurrentLimitVoltage = RegisterUInt16{Address: 2548, Sample: 40}
var RegisterLackVoltage = RegisterUInt16{Address: 2549, Sample: 40}
var RegisterOverVoltage = RegisterUInt16{Address: 2550, Sample: 88}
var RegisterRatedDCCurrent = RegisterUInt16{Address: 2551, Sample: 120}
var RegisterMaxDCCurrent = RegisterUInt16{Address: 2552, Sample: 120}
var RegisterMaxPhaseCurrent = RegisterUInt16{Address: 2553, Sample: 350}
var RegisterProtectedPhaseCurrent = RegisterUInt16{Address: 2554, Sample: 450}
var RegisterRatedPhaseCurrent = RegisterUInt16{Address: 2559, Sample: 150}
var RegisterControllerTemperature = RegisterFloat16{Address: 2562, Precision: Single, Sample: 32.97} // TODO: update with sample from your controller . There are two controller temperature -- which one is correct?
var RegisterUnderworkTemperature = RegisterUInt16{Address: 2568, Sample: 90}
var RegisterReworkTemperature = RegisterUInt16{Address: 2569, Sample: 80}
var RegisterLimitedCurrentTemperature = RegisterUInt16{Address: 2570, Sample: 70}
var RegisterMotorUnderworkTemperature = RegisterUInt16{Address: 2571, Sample: 130}
var RegisterMotorReworkTemperature = RegisterUInt16{Address: 2572, Sample: 110}
var RegisterMotorLimitedCurrentTemperature = RegisterUInt16{Address: 2573, Sample: 95}
var RegisterMotorTemperatureSensorEnabled = RegisterUInt16{Address: 2574, Sample: 0} // my controller doesn't have this register
var RegisterCurrentUnlockFlag = RegisterUInt16{Address: 2595, Sample: 3}
var RegisterReverseCurrentLimit = RegisterUInt16{Address: 2596, Sample: 20}
var RegisterFluxWeakeningEnabled = RegisterUInt16{Address: 2597, Sample: 1}
var RegisterFluxWeakenCurrent = RegisterUInt16{Address: 2598, Sample: 50}
var RegisterEbrakePhaseCurrent = RegisterUInt16{Address: 2599, Sample: 32}
var RegisterSlideRechargeEnabled = RegisterUInt16{Address: 2600, Sample: 0}
var RegisterSlideRechargePhaseCurrent = RegisterUInt16{Address: 2601, Sample: 30}
var RegisterSlideRechargeSpeed = RegisterUInt16{Address: 2602, Sample: 100}

var RegisterThrottleMinVoltage = RegisterFloat16{Address: 2608, Precision: Double, Sample: 0.75}
var RegisterThrottleMaxVoltage = RegisterFloat16{Address: 2609, Precision: Double, Sample: 4.0}
var RegisterAccelerateTime = RegisterUInt16{Address: 2610, Sample: 1}
var RegisterDecelerateTime = RegisterUInt16{Address: 2611, Sample: 1}

var RegisterThrottleMidVoltage = RegisterFloat16{Address: 2612, Precision: Double, Sample: 2.8}
var RegisterThrottleMidPhaseCurrent = RegisterUInt16{Address: 2613, Sample: 150}
var RegisterMotorPN = RegisterUInt16{Address: 2631, Sample: 16}
var RegisterMotorLmd = RegisterUInt16{Address: 2634, Sample: 50}
var RegisterSpeedLimitModeSelect = RegisterUInt16{Address: 2635, Sample: 0}
var RegisterMotorLimitSpeedSet = RegisterUInt16{Address: 2636, Sample: 100}
var RegisterLowSpeedSet = RegisterUInt16{Address: 2637, Sample: 70}
var RegisterMiddleSpeedSet = RegisterUInt16{Address: 2638, Sample: 100}
var RegisterCurrentLoopKp = RegisterUInt16{Address: 2648, Sample: 300}
var RegisterCurrentLoopKi = RegisterUInt16{Address: 2649, Sample: 9}
var RegisterHallAngleTestEnabled = RegisterUInt16{Address: 2650, Sample: 0}

var RegisterControlMode = RegisterUInt16{Address: 2651, Sample: 0} // 0=Normal 1=HallAngleTest
var RegisterTestCurrent = RegisterUInt16{Address: 2652, Sample: 0}
var RegisterHallAngle = RegisterUInt16{Address: 2653, Sample: 243}
var RegisterSystemStatus = RegisterUInt16{Address: 2748, Sample: 23}

/*
	0 PowerUpNoFinishedLackVoltage
	1 SystemError
	3 ThrottleProtect
	4 HallAngleTest
	5 CurrentDebug
	7 VoltDebug
	20 ElectronicBrake
	21 StopBrake
	22 SlideRecharge
	23 RunningFluxWeakenEnabled
	24 RunningFluxWeakenDisabled
	25 MotorReverse
	26 BrakeProtect
	27 GuardAgainstTheft
*/

var RegisterBatteryVoltage = RegisterFloat16{Address: 2749, Precision: Single, Sample: 67.63736264}
var RegisterWeakenCurrentCommand = RegisterUInt16{Address: 2750, Sample: 0}
var RegisterWeakenCurrentFeedback = RegisterUInt16{Address: 2751, Sample: 75}
var RegisterTorqueCurrentCommand = RegisterUInt16{Address: 2752, Sample: 0}
var RegisterTorqueCurrentFeedback = RegisterUInt16{Address: 2753, Sample: 65529}

// var RegisterControllerTemperature = RegisterUInt16{Address: 2754, Type: Int16, Sample: 25} // ** there are two controller temp -- which is correct according to the mqcon windows software?
var RegisterMotorTemperature = RegisterUInt16{Address: 2755, Sample: 130}
var RegisterMotorAngle = RegisterUInt16{Address: 2756, Sample: 22118}
var RegisterMotorSpeed = RegisterUInt16{Address: 2757, Sample: 0}
var RegisterHallStatus = RegisterUInt16{Address: 2758, Sample: 2}

var RegisterThrottleVoltage = RegisterFloat16{Address: 2759, Precision: Double, Sample: 0.128}
var RegisterMOSFETStatus = RegisterUInt16{Address: 2760, Sample: 426}
var RegisterInitial = RegisterUInt16{Address: 4039, Sample: 13345}

/* errors (where do we read this?)
1 MOSFETFault
2 OverVolt
3 LackVolt
5 MotorOverTemperature
6 ControllerOverTemperature
8 OverCurrentFault
9 Overload
11 StoreFault
12 HallTestFault
13 HallFault
18 OverSpeed
20 BlockProtectFault
21 UnInitEEPROM
25 PowerUpNoFinished
26 Brake
27 AntiTheft
28 Reverse
29 ReleaseThrottleError
30 ThrottleError
*/

/* inputs (where do we read this?

1 reverse
2 brake
3 boost
4 cruise
5 reverse
6 reverse
7 speed limit
8 antitheft
?? throttle?
*/

/* unknowns

2555 value: 51278
2556 value: 30 limit DC current?
2557 value: 150 ?? this is 100 on the spreadsheet; nominal battery current limit?
2558 value: 3000
2560 value: 30
2561 value: 4412 float
2588 value: 1
2589 value: 0
2590 value: 1
2591 value: 1
2592 value: 1
2593 value: 0
2594 value: 0
2603 value: 1
2614 value: 0
2628 value: 0
2629 value: 0
2630 value: 1
2632 value: 3072
2633 value: 102

2639 value: 0
2640 value: 2047

2654 value: 0
2655 value: 0
2656 value: 0
2657 value: 3
2658 value: 0
2659 value: 0
2660 value: 0
2661 value: 0
2662 value: 0

2761 value: 62970
2762 value: 0
2763 value: 0
2764 value: 0
2765 value: 0
2766 value: 0
2767 value: 65144
2768 value: 54452
2769 value: 50202
2770 value: 52916
2771 value: 50937
2772 value: 47101
2773 value: 0
2788 value: 0
2789 value: 14496
2790 value: 19560
2791 value: 0
2792 value: 32507
2793 value: 21819
2794 value: 0
2795 value: 22423
2796 value: 6627
2797 value: 0
2798 value: 0
2799 value: 0
2808 value: 0
2809 value: 53482
2810 value: 2000
2811 value: 45220
2812 value: 43306
2813 value: 44654
2814 value: 26965
2815 value: 3015
2816 value: 2736
2817 value: 54120
2848 value: 80
2849 value: 80
2850 value: 30
2851 value: 120
2852 value: 80
2853 value: 200
2854 value: 250
2855 value: 0
2856 value: 0
2857 value: 0
2858 value: 100
2868 value: 86
2869 value: 80
2870 value: 130
2871 value: 120
2872 value: 120
2873 value: 150
2874 value: 300
2875 value: 450
2876 value: 100
2877 value: 23
2878 value: 66
2879 value: 299
2880 value: 0
2881 value: 65
2882 value: 1
2883 value: 40
4038 value: 44039

4048 value: 0
*/
