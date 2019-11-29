package binutil

import (
	"fmt"
	"strings"
)

var ReverseOpMap map[int16]string
var OpMap map[string]int16
var OpMap2 map[string]int16
var OpList []string
var OpList2 []string

var OpList32  []string
var OpList64  []string
var OpListBr  []string
var OpListMem []string

var OpList_RA32ID   []string
var OpList_RB32ID   []string
var OpList_RP32ID   []string
var OpList_RM32ID   []string
var OpList_RA64ID   []string
var OpList_RB64ID   []string
var OpList_RC64ID   []string
var OpList_RM64ID   []string
var OpList_RD32ID   []string
var OpList_RD64ID   []string
var OpList_RD32ID_M []string
var OpList_RD64ID_M []string

func setConstMap(k string, v int16) {
	k=strings.TrimSpace(k)
	OpMap[k] = v
	kl := strings.ToLower(k)
	OpList = append(OpList, kl)
	OpMap[kl] = v
	ReverseOpMap[v] = kl
}
func setConstMap2(k string, v int16) {
	k=strings.TrimSpace(k)
	OpMap2[k] = v
	kl := strings.ToLower(k)
	OpList2 = append(OpList, kl)
	OpMap2[kl] = v
}

func init() {
	initOpMap()
	initOpMap2()
	initReverseUniOpMap()
	//classifyOpList()
}

func initOpMap() {
	OpMap = make(map[string]int16)
	ReverseOpMap = make(map[int16]string)
	setConstMap("BR_DS_IF0	",   0x000)
	setConstMap("BR_DS_IF1	",   0x001)
	setConstMap("BR_NDS_IF0	",   0x002)
	setConstMap("BR_NDS_IF1	",   0x003)
	setConstMap("JMP	",   0x004)
	setConstMap("CALL	",   0x005)
	setConstMap("AND_32	",   0x006)
	setConstMap("AND_64	",   0x007)
	setConstMap("OR_32	",   0x008)
	setConstMap("OR_64	",   0x009)
	setConstMap("XOR_32	",   0x00A)
	setConstMap("XOR_64	",   0x00B)
	setConstMap("UNISRC_32	",   0x00C)
	setConstMap("UNISRC_64	",   0x00D)
	setConstMap("SHL_32	",   0x00E)
	setConstMap("SHL_64	",   0x00F)
	setConstMap("SHR_S_32	",   0x010)
	setConstMap("SHR_S_64	",   0x011)
	setConstMap("SHR_U_32	",   0x012)
	setConstMap("SHR_U_64	",   0x013)
	setConstMap("ROTL_32	",   0x014)
	setConstMap("ROTL_64	",   0x015)
	setConstMap("ROTR_32	",   0x016)
	setConstMap("ROTR_64	",   0x017)
	setConstMap("SLT_S_32	",   0x018)
	setConstMap("SLT_S_64	",   0x019)
	setConstMap("SLT_U_32	",   0x01A)
	setConstMap("SLT_U_64	",   0x01B)
	setConstMap("SGE_S_32	",   0x01C)
	setConstMap("SGE_S_64	",   0x01D)
	setConstMap("SGE_U_32	",   0x01E)
	setConstMap("SGE_U_64	",   0x01F)
	setConstMap("SEQ_32	",   0x020)
	setConstMap("SEQ_64	",   0x021)
	setConstMap("SNE_32	",   0x022)
	setConstMap("SNE_64	",   0x023)
	setConstMap("BSLTI_S_32	",   0x024)
	setConstMap("BSLTI_S_64	",   0x025)
	setConstMap("BSLTI_U_32	",   0x026)
	setConstMap("BSLTI_U_64	",   0x027)
	setConstMap("BSGEI_S_32	",   0x028)
	setConstMap("BSGEI_S_64	",   0x029)
	setConstMap("BSGEI_U_32	",   0x02A)
	setConstMap("BSGEI_U_64	",   0x02B)
	setConstMap("BSEQI_32	",   0x02C)
	setConstMap("BSEQI_64	",   0x02D)
	setConstMap("BSNEI_32	",   0x02E)
	setConstMap("BSNEI_64	",   0x02F)
	setConstMap("ADDI_32	",   0x030)
	setConstMap("ADDI_64	",   0x031)
	setConstMap("SUBI_32	",   0x040)
	setConstMap("SUBI_64	",   0x041)
	setConstMap("ANDI_32	",   0x050)
	setConstMap("ANDI_64	",   0x051)
	setConstMap("ORI_32	",   0x060)
	setConstMap("ORI_64	",   0x061)
	setConstMap("XORI_32	",   0x070)
	setConstMap("XORI_64	",   0x071)
	setConstMap("SHLI_32	",   0x080)
	setConstMap("SHLI_64	",   0x081)
	setConstMap("SHRI_32	",   0x090)
	setConstMap("SHRI_64	",   0x091)
	setConstMap("ROTLI_32	",   0x0A0)
	setConstMap("ROTLI_64	",   0x0A1)
	setConstMap("SELECT_32	",   0x0B0)
	setConstMap("SELECT_64	",   0x0B1)
	setConstMap("CMOV_32	",   0x0B2)
	setConstMap("CMOV_64	",   0x0B3)
	setConstMap("MIN_S_32	",   0x0B4)
	setConstMap("MIN_S_64	",   0x0B5)
	setConstMap("MIN_U_32	",   0x0B6)
	setConstMap("MIN_U_64	",   0x0B7)
	setConstMap("MAX_S_32	",   0x0B8)
	setConstMap("MAX_S_64	",   0x0B9)
	setConstMap("MAX_U_32	",   0x0BA)
	setConstMap("MAX_U_64	",   0x0BB)
	setConstMap("ADD_32	",   0x0C0)
	setConstMap("ADD_64	",   0x0C1)
	setConstMap("SUB_32	",   0x0C2)
	setConstMap("SUB_64	",   0x0C3)
	setConstMap("ADC_32	",   0x0C4)
	setConstMap("ADC_64	",   0x0C5)
	setConstMap("SBB_32	",   0x0C6)
	setConstMap("SBB_64	",   0x0C7)
	setConstMap("MUL_32	",   0x0C8)
	setConstMap("MUL_64	",   0x0C9)
	setConstMap("LOAD_32	",   0x0D0)
	setConstMap("LOAD_64	",   0x0D1)
	setConstMap("F_LOAD_32	",   0x0D2)
	setConstMap("F_LOAD_64	",   0x0D3)
	setConstMap("STORE_32	",   0x0D4)
	setConstMap("STORE_64	",   0x0D5)
	setConstMap("F_STORE_32	",   0x0D6)
	setConstMap("F_STORE_64	",   0x0D7)
	setConstMap("LOADIMM_S_32",  0x0D8)
	setConstMap("LOADIMM_S_64",  0x0D9)
	setConstMap("LOADIMM_U_32",  0x0DA)
	setConstMap("LOADIMM_U_64",  0x0DB)
	setConstMap("LOAD8_S_32	",   0x0DC)
	setConstMap("LOAD8_S_64	",   0x0DD)
	setConstMap("LOAD8_U_32	",   0x0DE)
	setConstMap("LOAD8_U_64	",   0x0DF)
	setConstMap("LOAD16_S_32",   0x0E0)
	setConstMap("LOAD16_S_64",   0x0E1)
	setConstMap("LOAD16_U_32",   0x0E2)
	setConstMap("LOAD16_U_64",   0x0E3)
	setConstMap("STORE8_32	",   0x0E4)
	setConstMap("STORE8_64	",   0x0E5)
	setConstMap("STORE16_32	",   0x0E6)
	setConstMap("STORE16_64	",   0x0E7)
	setConstMap("STORE32_64	",   0x0E8)
	setConstMap("LOAD32_S_64",   0x0E9)
	setConstMap("LOAD32_U_64",   0x0EA)
	setConstMap("LOAD128_64	",   0x0EB)
	setConstMap("STORE128_64",   0x0EC)
	setConstMap("XMUL_S	",   0x0F0)
	setConstMap("XMUL_U	",   0x0F1)
	setConstMap("XMADD_S	",   0x0F0)
	setConstMap("XMADD_U	",   0x0F1)
	setConstMap("XMSUB_S	",   0x0F2)
	setConstMap("XMSUB_U	",   0x0F3)
	setConstMap("XNMADD_S	",   0x0F4)
	setConstMap("XNMADD_U	",   0x0F5)
	setConstMap("XNMSUB_S	",   0x0F6)
	setConstMap("XNMSUB_U	",   0x0F7)
	setConstMap("XMADC_S	",   0x0F8)
	setConstMap("XMADC_U	",   0x0F9)
	setConstMap("XMSBB_S	",   0x0FA)
	setConstMap("XMSBB_U	",   0x0FB)
	setConstMap("XNMADC_S	",   0x0FC)
	setConstMap("XNMADC_U	",   0x0FD)
	setConstMap("XNMSBB_S	",   0x0FE)
	setConstMap("XNMSBB_U	",   0x0FF)
	setConstMap("F.ADD_32	",   0x180)
	setConstMap("F.ADD_64	",   0x181)
	setConstMap("F.SUB_32	",   0x182)
	setConstMap("F.SUB_64	",   0x183)
	setConstMap("F.MUL_32	",   0x184)
	setConstMap("F.MUL_64	",   0x185)
	setConstMap("F.MIN_32	",   0x186)
	setConstMap("F.MIN_64	",   0x187)
	setConstMap("F.MAX_32	",   0x188)
	setConstMap("F.MAX_64	",   0x189)
	setConstMap("F.COPYSIGN_32", 0x18A)
	setConstMap("F.COPYSIGN_64", 0x18B)
	setConstMap("F.UNISRC_32",   0x18C)
	setConstMap("F.UNISRC_64",   0x18D)
	setConstMap("F.SEQ_32	",   0x190)
	setConstMap("F.SEQ_64	",   0x191)
	setConstMap("F.SNE_32	",   0x192)
	setConstMap("F.SNE_64	",   0x193)
	setConstMap("F.SLT_32	",   0x194)
	setConstMap("F.SLT_64	",   0x195)
	setConstMap("F.SGE_32	",   0x196)
	setConstMap("F.SGE_64	",   0x197)
	setConstMap("F.CVT_32	",   0x198)
	setConstMap("F.CVT_64	",   0x199)
	setConstMap("F.MADD_32	",   0x1A0)
	setConstMap("F.MADD_64	",   0x1A1)
	setConstMap("F.MSUB_32	",   0x1A2)
	setConstMap("F.MSUB_64	",   0x1A3)
	setConstMap("F.NMADD_32	",   0x1A8)
	setConstMap("F.NMADD_64	",   0x1A9)
	setConstMap("F.NMSUB_32	",   0x1AA)
	setConstMap("F.NMSUB_64	",   0x1AB)
}

func initOpMap2() {
	setConstMap2("YARNID	",  0)
	setConstMap2("NOT	",  1)
	setConstMap2("NEG	",  2)
	setConstMap2("CLZ	",  3)
	setConstMap2("CTZ	",  4)
	setConstMap2("POPCNT	",  5)
	setConstMap2("MOV64LOTO32", 6)
	setConstMap2("MOV32TO64_S", 6)
	setConstMap2("MOV64HITO32", 7)
	setConstMap2("MOV32TO64_U", 7)
	setConstMap2("CLEARBOOL	",  8)
	setConstMap2("SETBOOL	",  9)

	setConstMap2("F.FLOOR  	",  0)
	setConstMap2("F.CEIL  	",  1)
	setConstMap2("F.TRUNC   ",  2)
	setConstMap2("F.ROUND   ",  3)
	setConstMap2("F.ABS     ",  4)
	setConstMap2("F.NEG     ",  5)

	setConstMap2("FROMU32   ",  0)
	setConstMap2("FROMU64   ",  1)
	setConstMap2("FROMI32   ",  2)
	setConstMap2("FROMI64   ",  3)
	setConstMap2("TOU32     ",  4)
	setConstMap2("TOU64     ",  5)
	setConstMap2("TOI32     ",  6)
	setConstMap2("TOI64     ",  7)
	setConstMap2("TOANOTHER ",  8)
}

var UniOpMap = map[string][2]string{
	"yarnid_32":  [2]string{"UNISRC_32","YARNID"},
	"yarnid_64":  [2]string{"UNISRC_64","YARNID"},
	"not_32":  [2]string{"UNISRC_32","NOT"},
	"not_64":  [2]string{"UNISRC_64","NOT"},
	"neg_32":  [2]string{"UNISRC_32","NEG"},
	"neg_64":  [2]string{"UNISRC_64","NEG"},
	"clz_32":  [2]string{"UNISRC_32","CLZ"},
	"clz_64":  [2]string{"UNISRC_64","CLZ"},
	"ctz_32":  [2]string{"UNISRC_32","CTZ"},
	"ctz_64":  [2]string{"UNISRC_64","CTZ"},
	"popcnt_32":  [2]string{"UNISRC_32","POPCNT"},
	"popcnt_64":  [2]string{"UNISRC_64","POPCNT"},
	"mov64loto32":  [2]string{"UNISRC_32","MOV64LOTO32"},
	"mov64hito32":  [2]string{"UNISRC_32","MOV64HITO32"},
	"mov32to64_s":  [2]string{"UNISRC_64","MOV32TO64_S"},
	"mov32to64_u":  [2]string{"UNISRC_64","MOV32TO64_U"},
	"f.floor_32":  [2]string{"F.UNISRC_32","F.FLOOR"},
	"f.floor_64":  [2]string{"F.UNISRC_64","F.FLOOR"},
	"f.ceil_32":  [2]string{"F.UNISRC_32","F.CEIL"},
	"f.ceil_64":  [2]string{"F.UNISRC_64","F.CEIL"},
	"f.trunc_32":  [2]string{"F.UNISRC_32","F.TRUNC"},
	"f.trunc_64":  [2]string{"F.UNISRC_64","F.TRUNC"},
	"f.round_32":  [2]string{"F.UNISRC_32","F.ROUND"},
	"f.round_64":  [2]string{"F.UNISRC_64","F.ROUND"},
	"f.abs_32":  [2]string{"F.UNISRC_32","F.ABS"},
	"f.abs_64":  [2]string{"F.UNISRC_64","F.ABS"},
	"f.neg_32":  [2]string{"F.UNISRC_32","F.NEG"},
	"f.neg_64":  [2]string{"F.UNISRC_64","F.NEG"},
	"f.u32tof32":  [2]string{"F.CVT_32","FROMU32"},
	"f.u32tof64":  [2]string{"F.CVT_64","FROMU32"},
	"f.i32tof32":  [2]string{"F.CVT_32","FROMI32"},
	"f.i32tof64":  [2]string{"F.CVT_64","FROMI32"},
	"f.u64tof32":  [2]string{"F.CVT_32","FROMU64"},
	"f.u64tof64":  [2]string{"F.CVT_64","FROMU64"},
	"f.i64tof32":  [2]string{"F.CVT_32","FROMI64"},
	"f.i64tof64":  [2]string{"F.CVT_64","FROMI64"},
	"f.f32tou32":  [2]string{"F.CVT_32","TOU32"},
	"f.f64tou32":  [2]string{"F.CVT_64","TOU32"},
	"f.f32toi32":  [2]string{"F.CVT_32","TOI32"},
	"f.f64toi32":  [2]string{"F.CVT_64","TOI32"},
	"f.f32tou64":  [2]string{"F.CVT_32","TOU64"},
	"f.f64tou64":  [2]string{"F.CVT_64","TOU64"},
	"f.f32toi64":  [2]string{"F.CVT_32","TOI64"},
	"f.f64toi64":  [2]string{"F.CVT_64","TOI64"},
	"f.f32tof64":  [2]string{"F.CVT_32","TOANOTHER"},
	"f.f64tof32":  [2]string{"F.CVT_64","TOANOTHER"},
}

var ReverseUniOpMap map[int16]string

func initReverseUniOpMap() {
	for instr, twostr := range UniOpMap {
		op, op2 := OpMap[twostr[0]], OpMap2[twostr[1]]
		op = (op<<5) | op2
		ReverseUniOpMap[op] = instr
	}
}

func classifyOpList() {
	instr := Instr{
		RD:  1,
		RA:  1,
		RB:  1,
		Imm: 1,
		BrTarget: "not-empty",
	}
	for _, op := range OpList {
		dg := NewDecodedInstrGroup()
		instr.Mnemonic = op
		DecodeInstrForGroup(instr, &dg)
		if len(dg.InstrMem) != 0 {
			OpListMem = append(OpListMem, op)
		}
		if len(dg.InstrBr) != 0 {
			OpListBr = append(OpListBr, op)
		}
		if len(dg.Instr32) != 0 {
			OpList32 = append(OpList32, op)
		}
		if len(dg.Instr64) != 0 {
			OpList64 = append(OpList64, op)
		}
		if dg.RA32ID   != 0 { OpList_RA32ID   = append(OpList_RA32ID   , op) }
		if dg.RB32ID   != 0 { OpList_RB32ID   = append(OpList_RB32ID   , op) }
		if dg.RP32ID   != 0 { OpList_RP32ID   = append(OpList_RP32ID   , op) }
		if dg.RM32ID   != 0 { OpList_RM32ID   = append(OpList_RM32ID   , op) }
		if dg.RA64ID   != 0 { OpList_RA64ID   = append(OpList_RA64ID   , op) }
		if dg.RB64ID   != 0 { OpList_RB64ID   = append(OpList_RB64ID   , op) }
		if dg.RC64ID   != 0 { OpList_RC64ID   = append(OpList_RC64ID   , op) }
		if dg.RM64ID   != 0 { OpList_RM64ID   = append(OpList_RM64ID   , op) }
		if dg.RD32ID   != 0 { OpList_RD32ID   = append(OpList_RD32ID   , op) }
		if dg.RD64ID   != 0 { OpList_RD64ID   = append(OpList_RD64ID   , op) }
		if dg.RD32ID_M != 0 { OpList_RD32ID_M = append(OpList_RD32ID_M , op) }
		if dg.RD64ID_M != 0 { OpList_RD64ID_M = append(OpList_RD64ID_M , op) }
	}
}

func printUseFunctions() {
	printUseFunction("decUsedInMem", OpListMem);
	printUseFunction("decUsedInBr",  OpListBr);
	printUseFunction("decUsedIn64",  OpList64);
	printUseFunction("decUsedIn32",  OpList32);
	printUseFunction("decUseRA32",  OpList_RA32ID);
	printUseFunction("decUseRB32",  OpList_RB32ID);
	printUseFunction("decUseRP32",  OpList_RP32ID);
	printUseFunction("decUseRM32",  OpList_RM32ID);
	printUseFunction("decUseRA64",  OpList_RA64ID);
	printUseFunction("decUseRB64",  OpList_RB64ID);
	printUseFunction("decUseRC64",  OpList_RC64ID);
	printUseFunction("decUseRM64",  OpList_RM64ID);
	printUseFunction("decUseRD32",  OpList_RD32ID);
	printUseFunction("decUseRD64",  OpList_RD64ID);
	printUseFunction("decUseRD32_M",  OpList_RD32ID_M);
	printUseFunction("decUseRD64_M",  OpList_RD64ID_M);
}

func printUseFunction(funcName string, opList []string) {
	fmt.Printf("function automatic logic %s(input [8:0] op);\n", funcName)
	fmt.Printf("return \n")
	for _, op := range opList {
		fmt.Printf("  op == `%s || \n", strings.ToUpper(op))
	}
	fmt.Printf("1'b0; \n")
}

func printConstDefine() {
	for _, op := range OpList {
		fmt.Printf("`define %s 9'h%x\n", op, OpMap[op])
	}
	for _, op := range OpList2 {
		fmt.Printf("`define %s 5'h%x\n", op, OpMap2[op])
	}
}

func PrintVerilogHead() {
	printConstDefine()
	printUseFunctions()
}
