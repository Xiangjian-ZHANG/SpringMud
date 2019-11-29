package binutil

import (
	"fmt"
	"strings"
)


type DecodedInstrGroup struct {
	Label      string
	PC         int
	LineNumber int

	RA32ID int8
	RB32ID int8
	RP32ID int8
	RM32ID int8

	RA64ID int8
	RB64ID int8
	RC64ID int8
	RM64ID int8

	RD32ID   int8
	RD64ID   int8
	RD32ID_M int8
	RD64ID_M int8

	Imm16For32  int16
	Imm16For64  int16
	Imm16ForMem int16

	BrTarget string

	InstrMem string
	InstrBr  string
	Instr32  string
	Instr64  string
}

// ==== 64 ====
// phase0: wr-port-for-exe
// phase1: rd-port-A or rdwr-port-for-mem
// phase2: rd-port-B
// phase3: rd-port-C or rd-port-A
// ==== 32 ====
// phase0: wr-port-for-exe
// phase1: rdwr-port-for-mem
// phase2: rd-port-A or pointer-port
// phase3: rd-port-B

func (dg *DecodedInstrGroup) CheckRFPortConflict() {
	if dg.RC64ID !=-1 && dg.RM64ID >= 4 && dg.RA64ID >= 4 {
		fmt.Printf("Line %d: load/store register and mac-RA regsiter conflicts\n", dg.LineNumber)
		panic("RF port resource conflict")
	}
	if dg.RM32ID >= 4 && dg.RP32ID >= 4 {
		fmt.Printf("Line %d: pointer register and 32bit RA regsiter conflicts\n", dg.LineNumber)
		panic("RF port resource conflict")
	}
}

const InvalidImm16 = int16(-1)<<15

func NewDecodedInstrGroup() DecodedInstrGroup {
	return DecodedInstrGroup{
		RA32ID:   -1,
		RB32ID:   -1,
		RP32ID:   -1,
		RM32ID:   -1,
		RA64ID:   -1,
		RB64ID:   -1,
		RC64ID:   -1,
		RM64ID:   -1,
		RD32ID:   -1,
		RD64ID:   -1,
		RD32ID_M: -1,
		RD64ID_M: -1,

		Imm16For32:  InvalidImm16,
		Imm16For64:  InvalidImm16,
		Imm16ForMem: InvalidImm16,
	}
}

func (dg *DecodedInstrGroup) SetRA32ID(id int8) {dg.CheckRegID(dg.RA32ID, "RA32ID"); dg.RA32ID=id}
func (dg *DecodedInstrGroup) SetRB32ID(id int8) {dg.CheckRegID(dg.RB32ID, "RB32ID"); dg.RB32ID=id}
func (dg *DecodedInstrGroup) SetRP32ID(id int8) {dg.CheckRegID(dg.RP32ID, "RP32ID"); dg.RP32ID=id}
func (dg *DecodedInstrGroup) SetRM32ID(id int8) {dg.CheckRegID(dg.RM32ID, "RC32ID"); dg.RM32ID=id}
func (dg *DecodedInstrGroup) SetRA64ID(id int8) {dg.CheckRegID(dg.RA64ID, "RA64ID"); dg.RA64ID=id}
func (dg *DecodedInstrGroup) SetRB64ID(id int8) {dg.CheckRegID(dg.RB64ID, "RB64ID"); dg.RB64ID=id}
func (dg *DecodedInstrGroup) SetRC64ID(id int8) {dg.CheckRegID(dg.RC64ID, "RC64ID"); dg.RC64ID=id}
func (dg *DecodedInstrGroup) SetRM64ID(id int8) {dg.CheckRegID(dg.RM64ID, "RM64ID"); dg.RM64ID=id}
func (dg *DecodedInstrGroup) SetRD32ID(id int8) {dg.CheckRegID(dg.RD32ID, "RD32ID"); dg.RD32ID=id}
func (dg *DecodedInstrGroup) SetRD64ID(id int8) {dg.CheckRegID(dg.RD64ID, "RD64ID"); dg.RD64ID=id}
func (dg *DecodedInstrGroup) SetRD32ID_M(id int8) {dg.CheckRegID(dg.RD32ID_M, "RD32ID_M"); dg.RD32ID_M=id}
func (dg *DecodedInstrGroup) SetRD64ID_M(id int8) {dg.CheckRegID(dg.RD64ID_M, "RD64ID_M"); dg.RD64ID_M=id}

func (dg *DecodedInstrGroup) SetInstrMem(instr string) {
	if dg.InstrMem != "" {
		fmt.Printf("Line %d: Instruction-for-Mem has been set. Resource conflict.\n", dg.LineNumber)
		panic("VLIW Execution-Unit Resource conflict")
	}
	dg.InstrMem = instr
}
func (dg *DecodedInstrGroup) SetInstrBr(instr string) {
	if dg.InstrBr != "" {
		fmt.Printf("Line %d: Instruction-for-Br has been set. Resource conflict.\n", dg.LineNumber)
		panic("VLIW Execution-Unit Resource conflict")
	}
	dg.InstrBr = instr
}
func (dg *DecodedInstrGroup) SetInstr32(instr string) {
	if dg.Instr32 != "" {
		fmt.Printf("Line %d: Instruction-for-32 has been set. Resource conflict.\n", dg.LineNumber)
		panic("VLIW Execution-Unit Resource conflict")
	}
	dg.Instr32 = instr
}
func (dg *DecodedInstrGroup) SetInstr64(instr string) {
	if dg.Instr64 != "" {
		fmt.Printf("Line %d: Instruction-for-64 has been set. Resource conflict.\n", dg.LineNumber)
		panic("VLIW Execution-Unit Resource conflict")
	}
	dg.Instr64 = instr
}

func (dg *DecodedInstrGroup) SetImm16ForMem(imm int16) {
	if dg.Imm16ForMem != InvalidImm16 {
		fmt.Printf("Line %d: Immediate-value-for-Mem has been set. Resource conflict.\n", dg.LineNumber)
		panic("VLIW Resource conflict")
	}
	dg.Imm16ForMem = imm
}
func (dg *DecodedInstrGroup) SetImm16For32(imm int16) {
	if dg.Imm16For32 != InvalidImm16 {
		fmt.Printf("Line %d: Immediate-value-for-32 has been set. Resource conflict.\n", dg.LineNumber)
		panic("VLIW Resource conflict")
	}
	dg.Imm16For32 = imm
}
func (dg *DecodedInstrGroup) SetImm16For64(imm int16) {
	if dg.Imm16For64 != InvalidImm16 {
		fmt.Printf("Line %d: Immediate-value-for-64 has been set. Resource conflict.\n", dg.LineNumber)
		panic("VLIW Resource conflict")
	}
	dg.Imm16For64 = imm
}

func (dg *DecodedInstrGroup) SetBrTarget(target string) {
	if len(dg.BrTarget) != 0 {
		fmt.Printf("Line %d: Branch-target has been set. Resource conflict.\n", dg.BrTarget)
		panic("VLIW Resource conflict")
	}
	dg.BrTarget = target
}

func (dg *DecodedInstrGroup) CheckRegID(id int8, name string) {
	if id!=-1 {
		fmt.Printf("Line %d: '%s' has been set. Resource conflict.\n", dg.LineNumber, name)
		panic("VLIW RegFile-Port Resource conflict")
	}
}

type InstrType int

const (
	UnknownInstr InstrType = iota
	Branch
	WrBoolReg
	RBImm10
	RDRBImm8
	RDRARB
	FP_RDRARB
	MAC
	Load
	Store
	YarnID
	RDRB
	FP_RDRB
	From64To32
	From32To64
)

func GetInstrType(mnemonic string) InstrType {
	switch mnemonic {
	case
	"br_ds_if0",
	"br_ds_if1",
	"br_nds_if0",
	"br_nds_if1",
	"jmp",
	"ret",
	"call":
		return Branch
	case
	"clearbool",
	"setbool":
		return WrBoolReg
	case
	"bslti_s_32",
	"bslti_s_64",
	"bslti_u_32",
	"bslti_u_64",
	"bsgei_s_32",
	"bsgei_s_64",
	"bsgei_u_32",
	"bsgei_u_64",
	"bseqi_32",
	"bseqi_64",
	"bsnei_32",
	"bsnei_64":
		return RBImm10
	case
	"addi_32",
	"addi_64",
	"subi_32",
	"subi_64",
	"andi_32",
	"andi_64",
	"ori_32",
	"ori_64",
	"xori_32",
	"xori_64",
	"shli_32",
	"shli_64",
	"shri_32",
	"shri_64",
	"rotli_32",
	"rotli_64":
		return RDRBImm8
	case
	"and_32",
	"and_64",
	"or_32",
	"or_64",
	"xor_32",
	"xor_64",
	"shl_32",
	"shl_64",
	"shr_s_32",
	"shr_s_64",
	"shr_u_32",
	"shr_u_64",
	"rotl_32",
	"rotl_64",
	"rotr_32",
	"rotr_64",
	"slt_s_32",
	"slt_s_64",
	"slt_u_32",
	"slt_u_64",
	"sge_s_32",
	"sge_s_64",
	"sge_u_32",
	"sge_u_64",
	"seq_32",
	"seq_64",
	"sne_32",
	"sne_64",
	"select_32",
	"select_64",
	"cmov_32",
	"cmov_64",
	"min_s_32",
	"min_s_64",
	"min_u_32",
	"min_u_64",
	"max_s_32",
	"max_s_64",
	"max_u_32",
	"max_u_64",
	"add_32",
	"add_64",
	"sub_32",
	"sub_64",
	"adc_32",
	"adc_64",
	"sbb_32",
	"sbb_64",
	"mul_32",
	"mul_64":
		return RDRARB
	case
	"f.add_32",
	"f.add_64",
	"f.sub_32",
	"f.sub_64",
	"f.mul_32",
	"f.mul_64",
	"f.min_32",
	"f.min_64",
	"f.max_32",
	"f.max_64",
	"f.copysign_32",
	"f.copysign_64",
	"f.seq_32",
	"f.seq_64",
	"f.sne_32",
	"f.sne_64",
	"f.slt_32",
	"f.slt_64",
	"f.sge_32",
	"f.sge_64":
		return FP_RDRARB
	case
	"f.madd_32",
	"f.madd_64",
	"f.msub_32",
	"f.msub_64",
	"f.nmadd_32",
	"f.nmadd_64",
	"f.nmsub_32",
	"f.nmsub_64",
	"xmul_s",
	"xmul_u",
	"xmadd_s",
	"xmadd_u",
	"xmsub_s",
	"xmsub_u",
	"xnmadd_s",
	"xnmadd_u",
	"xnmsub_s",
	"xnmsub_u",
	"xmadc_s",
	"xmadc_u",
	"xmsbb_s",
	"xmsbb_u",
	"xnmadc_s",
	"xnmadc_u",
	"xnmsbb_s",
	"xnmsbb_u":
		return MAC
	case
	"load_32",
	"load_64",
	"f.load_32",
	"f.load_64",
	"loadimm_s_32",
	"loadimm_s_64",
	"loadimm_u_32",
	"loadimm_u_64",
	"load8_s_32",
	"load8_s_64",
	"load8_u_32",
	"load8_u_64",
	"load16_s_32",
	"load16_s_64",
	"load16_u_32",
	"load16_u_64",
	"load32_s_32",
	"load32_u_64",
	"load128":
		return Load
	case
	"store128",
	"store_32",
	"store_64",
	"f.store_32",
	"f.store_64":
		return Store
	case
	"yarnid_32",
	"yarnid_64":
		return YarnID
	case
	"not_32",
	"not_64",
	"neg_32",
	"neg_64",
	"clz_32",
	"clz_64",
	"ctz_32",
	"ctz_64",
	"popcnt_32",
	"popcnt_64":
		return RDRB
	case
	"f.floor_32",
	"f.floor_64",
	"f.ceil_32",
	"f.ceil_64",
	"f.trunc_32",
	"f.trunc_64",
	"f.round_32",
	"f.round_64",
	"f.abs_32",
	"f.abs_64",
	"f.neg_32",
	"f.neg_64",
	"f.u64tof32",
	"f.u64tof64",
	"f.i64tof32",
	"f.i64tof64",
	"f.f32tou64",
	"f.f64tou64",
	"f.f32toi64",
	"f.f64toi64":
		return FP_RDRB
	case
	"f.f32tou32",
	"f.f64tou32",
	"f.f32toi32",
	"f.f64toi32",
	"f.u32tof32",
	"mov64loto32",
	"mov64hito32":
		return From64To32
	case
	"f.u32tof64",
	"f.i32tof32",
	"f.i32tof64",
	"f.f32tof64",
	"f.f64tof32",
	"mov32to64_s",
	"mov32to64_u":
		return From32To64
	default:
		return UnknownInstr
	}
}

func DecodeInstrForGroup(instr Instr, dg *DecodedInstrGroup) {
	switch GetInstrType(instr.Mnemonic) {
	case Branch:
		dg.SetBrTarget(instr.BrTarget)
		dg.SetInstrBr(instr.Mnemonic)
	case WrBoolReg:
		dg.SetInstr32(instr.Mnemonic)
	case RBImm10:
		if strings.HasSuffix(instr.Mnemonic, "_32") {
			dg.SetRB32ID(instr.RB)
			dg.SetImm16For32(instr.Imm)
			dg.SetInstr32(instr.Mnemonic)
		} else {
			dg.SetRB64ID(instr.RB)
			dg.SetImm16For64(instr.Imm)
			dg.SetInstr64(instr.Mnemonic)
		}
	case RDRBImm8:
		if strings.HasSuffix(instr.Mnemonic, "_32") {
			dg.SetRB32ID(instr.RB)
			dg.SetRD32ID(instr.RD)
			dg.SetImm16For32(instr.Imm)
			dg.SetInstr32(instr.Mnemonic)
		} else {
			dg.SetRB64ID(instr.RB)
			dg.SetRD64ID(instr.RD)
			dg.SetImm16For64(instr.Imm)
			dg.SetInstr64(instr.Mnemonic)
		}
	case RDRARB:
		if strings.HasSuffix(instr.Mnemonic, "_32") {
			dg.SetRA32ID(instr.RA)
			dg.SetRB32ID(instr.RB)
			dg.SetRD32ID(instr.RD)
			dg.SetInstr32(instr.Mnemonic)
		} else {
			dg.SetRA64ID(instr.RA)
			dg.SetRB64ID(instr.RB)
			dg.SetRD64ID(instr.RD)
			dg.SetInstr64(instr.Mnemonic)
		}
	case FP_RDRARB:
		dg.SetRA64ID(instr.RA)
		dg.SetRB64ID(instr.RB)
		dg.SetRD64ID(instr.RD)
		dg.SetInstr64(instr.Mnemonic)
	case MAC:
		dg.SetRA64ID(instr.RA)
		dg.SetRB64ID(instr.RB)
		dg.SetRD64ID(instr.RD)
		dg.SetRC64ID(instr.RD)
		dg.SetInstr64(instr.Mnemonic)
	case Load:
		if strings.HasSuffix(instr.Mnemonic, "_32") {
			dg.SetRD32ID_M(instr.RD)
		} else {
			dg.SetRD64ID_M(instr.RD)
		}
		dg.SetRP32ID(instr.RB)
		dg.SetImm16ForMem(instr.Imm)
		dg.SetInstrMem(instr.Mnemonic)
	case Store:
		if strings.HasSuffix(instr.Mnemonic, "_32") {
			dg.SetRM32ID(instr.RD)
		} else {
			dg.SetRM64ID(instr.RD)
		}
		dg.SetRP32ID(instr.RB)
		dg.SetImm16ForMem(instr.Imm)
		dg.SetInstrMem(instr.Mnemonic)
	case YarnID:
		if strings.HasSuffix(instr.Mnemonic, "_32") {
			dg.SetRD32ID(instr.RD)
		} else {
			dg.SetRD64ID(instr.RD)
		}
	case RDRB:
		if strings.HasSuffix(instr.Mnemonic, "_32") {
			dg.SetRB32ID(instr.RB)
			dg.SetRD32ID(instr.RD)
			dg.SetInstr32(instr.Mnemonic)
		} else {
			dg.SetRB64ID(instr.RB)
			dg.SetRD64ID(instr.RD)
			dg.SetInstr64(instr.Mnemonic)
		}
	case FP_RDRB:
		dg.SetRB64ID(instr.RB)
		dg.SetRD64ID(instr.RD)
		dg.SetInstr64(instr.Mnemonic)
	case From64To32:
		dg.SetRB64ID(instr.RB)
		dg.SetRD32ID(instr.RD)
		dg.SetInstr64(instr.Mnemonic)
	case From32To64:
		dg.SetRB32ID(instr.RB)
		dg.SetRD64ID(instr.RD)
		dg.SetInstr64(instr.Mnemonic)
	default:
		panic("Unknown Instruction")
	}
}

func DecodeInstrGroup(instrGroup InstrGroup) DecodedInstrGroup {
	dg := NewDecodedInstrGroup()
	dg.Label = instrGroup.Label
	dg.PC = instrGroup.PC
	dg.LineNumber = instrGroup.LineNumber
	for _, instr := range instrGroup.InstrList {
		DecodeInstrForGroup(instr, &dg)
	}
	return dg
}

func MnemonicToFuncName(mnemonic string) string {
	mnemonic = strings.ReplaceAll(mnemonic, ".", "_")
	return strings.ToUpper(mnemonic)+"()"
}

func (dg *DecodedInstrGroup) WriteSimCode(strb *strings.Builder, delayedBrTarget *DecodedInstrGroup, brTarget *DecodedInstrGroup) {
	var w = func(s string) {
		strb.WriteString(s)
	}
	if len(dg.Label) != 0 {
		w(dg.Label+":\n")
	}
	w(fmt.Sprintf("/* Source Line Number: %d */\n", dg.LineNumber))
	w(fmt.Sprintf("PC = %d;\n", dg.PC))
	w(CLEAR_WR_EN)
	w(fmt.Sprintf("RA32ID=%d; RB32ID=%d; RP32ID=%d; RM32ID=%d;\n", dg.RA32ID, dg.RB32ID, dg.RP32ID, dg.RM32ID))
	w(fmt.Sprintf("RA64ID=%d; RB64ID=%d; RC64ID=%d; RM64ID=%d;\n", dg.RA64ID, dg.RB64ID, dg.RC64ID, dg.RM64ID))
	w(fmt.Sprintf("RD32ID=%d; RD64ID=%d; RD32ID_M=%d; RD64ID_M=%d;\n", dg.RD32ID, dg.RD64ID, dg.RD32ID_M, dg.RD64ID_M))
	w(fmt.Sprintf("Imm16For32=%d; Imm16For64=%d; Imm16ForMem=%d;\n", dg.Imm16For32, dg.Imm16For64, dg.Imm16ForMem))
	w("for(YI=0; YI<YC; YI++) {\n")
	w("    if(PC<FwdBrTarget[YI]) continue; /*During Instruction Skipping*/\n")
	w("    if(PC==FwdBrTarget[YI]) FwdBrTarget[YI]=-1;/*End of Instruction Skipping*/\n")
	for _, mnemonic := range []string{dg.InstrMem, dg.Instr64, dg.Instr32} {
		if len(mnemonic)!=0 {
			w(MnemonicToFuncName(mnemonic))
		}
	}
	// BR_DS_IF0_PREPARE
	// BR_DS_IF1_PREPARE
	// BACK_BR_DS_EXECUTE
	// BACK_BR_NDS_IF0
	// BACK_BR_NDS_IF1
	// FWD_BR_NDS_IF0
	// FWD_BR_NDS_IF1
	// JMP CALL RET
	if delayedBrTarget != nil { // in a delay slot
		if len(dg.InstrBr)!=0 {
			fmt.Printf("Line %d: Cannot have control flow instructions in delay slot!\n", dg.LineNumber)
			panic("Branch in delay slot!")
		}
		if delayedBrTarget.PC > dg.PC { // delayed branch forward
			panic("Forwrd Branch with delay slot is not permitted!")
		}
	} else if len(dg.InstrBr) != 0 && brTarget.PC > dg.PC { // normal branch forward
		switch dg.InstrBr {
		case "br_ds_if0":
			panic("Forwrd Branch with delay slot is not permitted!")
		case "br_ds_if1":
			panic("Forwrd Branch with delay slot is not permitted!")
		case "br_nds_if0":
			w("FWD_BR_NDS_IF0("+brTarget.Label+")")
		case "br_nds_if1":
			w("FWD_BR_NDS_IF1("+brTarget.Label+")")
		}
	}

	w("}\n") // end of the loop over yarns
	w(DO_WRITE)

	if delayedBrTarget != nil { // in a delay slot 
		w("BACK_BR_DS_EXECUTE()") // forward branch with delay slot is not permitted
	}

	if dg.InstrBr == "br_ds_if0" && brTarget.PC < dg.PC {
		w("BACK_BR_DS_IF0_PREPARE()")
	} else if dg.InstrBr == "br_ds_if1" && brTarget.PC < dg.PC {
		w("BACK_BR_DS_IF1_PREPARE()")
	} else if dg.InstrBr == "br_nds_if0" && brTarget.PC < dg.PC {
		w("BACK_BR_NDS_IF0()")
	} else if dg.InstrBr == "br_nds_if1" && brTarget.PC < dg.PC {
		w("BACK_BR_NDS_IF1()")
	} else if dg.InstrBr == "jmp" {
		w("JMP("+dg.BrTarget+")")
	} else if dg.InstrBr == "ret" {
		w("RET()")
	} else if dg.InstrBr == "call" {
		w(fmt.Sprintf("CALL(%s,ENDOF_%d)\nENDOF_%d:", dg.BrTarget, dg.PC, dg.PC))
	}
}
