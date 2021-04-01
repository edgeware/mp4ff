package hevc

import (
	"fmt"
)

//Codecs - MIME subtype like hev1.1.6.L93.B0 where hev1 is sampleEntry string.
//
// Following ISO/IEC 14496-15 2017 Annex E
func Codecs(sampleEntry string, sps *SPS) string {
	profilePart := ""
	ptl := sps.ProfileTierLevel
	switch ptl.GeneralProfileSpace {
	case 0:
		// Nothing
	case 1:
		profilePart += "A"
	case 2:
		profilePart += "B"
	case 3:
		profilePart += "C"
	}
	profilePart += fmt.Sprintf("%d", ptl.GeneralProfileIDC)

	flagsPart := fmt.Sprintf("%X", reverseFlags(ptl.GeneralProfileCompatibilityFlags))
	var levelPart string
	if ptl.GeneralTierFlag {
		levelPart = "H"
	} else {
		levelPart = "L"
	}
	levelPart += fmt.Sprintf("%d", ptl.GeneralLevelIDC)
	cif := ptl.GeneralConstraintIndicatorFlags
	nrBytes := 6
	for i := 0; i < 5; i++ { // Remove trailing zero bytes
		if cif&0xff == 0 {
			cif = cif >> 8
			nrBytes--
		} else {
			break
		}
	}
	// Code remaining nrBytes as hex separated by dots
	constraintBytes := ""
	for i := 0; i < nrBytes; i++ {
		constraintBytes += fmt.Sprintf(".%X", (cif>>((nrBytes-1-i)*8))&0xff)
	}

	return fmt.Sprintf("%s.%s.%s.%s%s", sampleEntry, profilePart, flagsPart, levelPart, constraintBytes)
}

func reverseFlags(flags uint32) uint32 {
	var rf uint32 = 0
	for i := 0; i < 32; i++ {
		if flags&(1<<(31-i)) != 0 {
			rf |= 1 << i
		}
	}
	return rf
}
