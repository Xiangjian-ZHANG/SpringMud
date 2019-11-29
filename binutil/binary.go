package binutil

import (
	"fmt"
	"os"
	"io"
	"encoding/binary"
)

type RawInstr struct {
	OP       int16
	Imm      int16
	RD       int8
	RA       int8
	RB       int8
}

func DecodeRaw(bits uint32) (res RawInstr) {
	bits = bits >> 1 //Remove the group-ending bit
	res.Imm = int16(bits & 0x7FFF)
	res.RB = int8(bits & 0x1F)
	bits = bits >> 5
	res.RA = int8(bits & 0x1F)
	bits = bits >> 5
	res.RD = int8(bits & 0x1F)
	bits = bits >> 5
	res.OP = int16(bits & 0x1FF)
	return
}

// -------------- | ---------------- | -------- | -------- | ------ |
// Calc-Reg-Reg   | 9bit             | RD       | RA       | RB     |
// Calc-Reg       | 9bit             | RD       | extra op | RB     |
// Calc-Reg-Imm8  | 6bit op + imm7~5 | RD       | imm4~0   | RB     |
// Load           | 9bit             | RD       | imm4~0   | RB     |
// Store          | 9bit             | RM       | imm4~0   | RB     |
// Bool-Set-on-XX | 9bit             | imm9~5   | imm4~0   | RB     |
// ControlFlow    | 9bit             | imm14~10 | imm4~0   | imm9~5 |

func EncodeCalcRegImm8(instr Instr) (res uint32) {
	res = uint32(OpMap[instr.Mnemonic]) | uint32((instr.Imm&0b11100000)>>4)
	res = (res << 5) | uint32(instr.RD & 0x1F)
	res = (res << 5) | uint32(instr.Imm & 0x1F)
	res = (res << 5) | uint32(instr.RB & 0x1F)
	return
}

func DecodeCalcRegImm8(rawI RawInstr, mnemonic string) (res Instr) {
	res.Mnemonic = mnemonic
	res.RD = rawI.RD
	res.Imm = int16(((rawI.OP>>1)&0x7)<<5) | int16(rawI.RA)
	res.RB = rawI.RB
	return
}

func EncodeCalcRegReg(instr Instr) (res uint32) {
	res = uint32(OpMap[instr.Mnemonic])
	res = (res << 5) | uint32(instr.RD & 0x1F)
	res = (res << 5) | uint32(instr.RA & 0x1F)
	res = (res << 5) | uint32(instr.RB & 0x1F)
	return
}

func DecodeCalcRegReg(rawI RawInstr, mnemonic string) (res Instr) {
	res.Mnemonic = mnemonic
	res.RD = rawI.RD
	res.RA = rawI.RA
	res.RB = rawI.RB
	return
}

func EncodeCalcReg(instr Instr) (res uint32) {
	twostr, ok := UniOpMap[instr.Mnemonic]
	if !ok {
		panic("invalid instruction")
	}
	op, op2 := OpMap[twostr[0]], OpMap2[twostr[1]]
	res = uint32(op)
	res = (res << 5) | uint32(instr.RD & 0x1F)
	res = (res << 5) | uint32(op2 & 0x1F)
	res = (res << 5) | uint32(instr.RB & 0x1F)
	return
}

func DecodeCalcReg(rawI RawInstr, mnemonic string) (res Instr) {
	res.Mnemonic = mnemonic
	res.RD = rawI.RD
	res.RB = rawI.RB
	return
}

func EncodeLoadStore(instr Instr) (res uint32) {
	res = uint32(OpMap[instr.Mnemonic])
	res = (res << 5) | uint32(instr.RD & 0x1F)
	res = (res << 5) | uint32(instr.Imm & 0x1F)
	res = (res << 5) | uint32(instr.RB & 0x1F)
	return
}

func DecodeLoadStore(rawI RawInstr, mnemonic string) (res Instr) {
	res.Mnemonic = mnemonic
	res.Imm = int16(rawI.RA)
	res.Imm = (res.Imm<<11)>>11 // Sign extension
	res.RB = rawI.RB
	res.RD = rawI.RD
	return
}

func EncodeBoolSetOnXX(instr Instr) (res uint32) {
	res = uint32(OpMap[instr.Mnemonic])
	res = (res << 10) | uint32(instr.Imm & 0x3FF)
	res = (res << 5)  | uint32(instr.RB & 0x1F)
	return
}

func DecodeBoolSetOnXX(rawI RawInstr, mnemonic string) (res Instr) {
	res.Mnemonic = mnemonic
	res.RB = rawI.RB
	res.Imm = rawI.Imm >> 5
	return
}

func EncodeControlFlow(instr Instr) (res uint32) {
	res = uint32(OpMap[instr.Mnemonic])
	res = (res << 15) | uint32(instr.Imm & 0x7FFF)
	return
}

func DecodeControlFlow(rawI RawInstr, mnemonic string) (res Instr) {
	res.Mnemonic = mnemonic
	res.Imm = rawI.Imm
	if res.Mnemonic == "call" && res.Imm == 0 {
		res.Mnemonic = "ret"
	}
	return
}

func ToBinary(instrGroupList []InstrGroup) (res []uint32) {
	for _, instrGroup := range instrGroupList {
		lastIdx := len(instrGroup.InstrList)-1
		for i:=0; i<=lastIdx; i++ {
			bits := EncodeInstr(instrGroup.InstrList[i])
			bits = bits << 1;
			if i==lastIdx {
				bits = bits | 1 // add group-ending bit
			}
			res = append(res, bits)
		}
	}
	return
}

func FromBinary(binData []uint32) (instrGroupList []InstrGroup) {
	currInstrGroup := InstrGroup{}
	for _, bits := range binData {
		rawI := DecodeRaw(bits)
		mnemonic, ok := ReverseOpMap[rawI.OP]
		if !ok {
			fmt.Printf("Unknown OP code, bits: %x op: %x\n", bits, rawI.OP)
			panic("Invalid OP code")
		}
		if mnemonic == "UNISRC_32" || mnemonic == "UNISRC_64" ||
		mnemonic == "F.UNISRC_32" || mnemonic == "F.UNISRC_64" ||
		mnemonic == "F.CVT_32" || mnemonic == "F.CVT_64" {
			mnemonic, ok = ReverseUniOpMap[(rawI.OP<<5)|int16(rawI.RA)]
		}
		if !ok {
			fmt.Printf("Unknown OP code, bits: %x op: %x op2: %x\n", bits, rawI.OP, rawI.RA)
			panic("Invalid OP code")
		}
		instr := DecodeInstr(rawI, mnemonic)
		currInstrGroup.InstrList = append(currInstrGroup.InstrList, instr)
		if (bits & 1) != 0 {
			instrGroupList = append(instrGroupList, currInstrGroup)
			currInstrGroup = InstrGroup{}
		}
	}
	AssignPC(instrGroupList)
	AssignBrTarget(instrGroupList)
	return
}

func reportErrForInvalidTarget(instrGroup InstrGroup) {
	fmt.Printf("PC: %d, instructions:\n", instrGroup.PC)
	for _, instr := range instrGroup.InstrList {
		fmt.Println(FormatInstr(instr))
	}
}

func AssignBrTarget(instrGroupList []InstrGroup) {
	pc2group := make(map[int]*InstrGroup)
	for _, instrGroup := range instrGroupList {
		pc2group[instrGroup.PC] = &instrGroup
	}
	for _, instrGroup := range instrGroupList {
		for _, instr := range instrGroup.InstrList {
			if GetFormatType(instr.Mnemonic) != ControlFlow {
				continue
			}
			target := instrGroup.PC + int(instr.Imm)
			instr.BrTarget = fmt.Sprintf("fmt_%d", target)
			targetInstrGroup, ok := pc2group[target]
			if !ok {
				fmt.Println("Invalid Branch Target!")
				reportErrForInvalidTarget(instrGroup)
				panic("Invalid Branch Target")
			}
			targetInstrGroup.Label = fmt.Sprintf("pc_%d", targetInstrGroup.PC)
		}
	}
}

func DecodeInstr(rawI RawInstr, mnemonic string) (res Instr) {
	switch GetFormatType(mnemonic) {
	case ControlFlow:
		return DecodeControlFlow(rawI, mnemonic)
	case NoParam:
		if mnemonic=="ret" {
			return DecodeControlFlow(rawI, mnemonic)
		} else {
			return DecodeCalcReg(rawI, mnemonic)
		}
	case SrcRegImm10:
		return DecodeBoolSetOnXX(rawI, mnemonic)
	case SrcRegImm8:
		return DecodeCalcRegImm8(rawI, mnemonic)
	case SrcRegX2:
		return DecodeCalcRegReg(rawI, mnemonic)
	case LoadRegIMM5, StoreRegIMM5:
		return DecodeLoadStore(rawI, mnemonic)
	case OnlyDst, SrcRegX1:
		return DecodeCalcReg(rawI, mnemonic)
	default:
		panic("Unknown FormatType")
	}
}

func EncodeInstr(instr Instr) uint32 {
	switch GetFormatType(instr.Mnemonic) {
	case ControlFlow:
		return EncodeControlFlow(instr)
	case NoParam:
		if instr.Mnemonic=="ret" {
			return EncodeControlFlow(instr)
		} else {
			return EncodeCalcReg(instr)
		}
	case SrcRegImm10:
		return EncodeBoolSetOnXX(instr)
	case SrcRegImm8:
		return EncodeCalcRegImm8(instr)
	case SrcRegX2:
		return EncodeCalcRegReg(instr)
	case LoadRegIMM5, StoreRegIMM5:
		return EncodeLoadStore(instr)
	case OnlyDst,SrcRegX1:
		return EncodeCalcReg(instr)
	default:
		panic("Unknown FormatType")
	}
}

func ReadBinaryFile(fname string) (res []uint32) {
	f, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var buf [4]byte
	for {
		n, err := f.Read(buf[:])
		if n == 0 && err == io.EOF {
			break
		} else if n != 0 && err == io.EOF {
			panic("Unexpected End-of-File")
		} else if err != nil {
			panic(err)
		}
		res = append(res, binary.LittleEndian.Uint32(buf[:]))
	}
	return
}

func WriteBinaryFile(fname string, data []uint32) {
	f, err := os.OpenFile(fname, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var buf [4]byte
	for _, bits := range data {
		binary.LittleEndian.PutUint32(buf[:], bits)
		_, err := f.Write(buf[:])
		if err != nil {
			panic(err)
		}
	}
}
