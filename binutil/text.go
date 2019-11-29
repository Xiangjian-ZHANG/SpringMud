package binutil

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Instr struct {
	Mnemonic string
	RD       int8
	RA       int8
	RB       int8
	Imm      int16
	BrTarget string
}

const (
	BeginGroup = 0
	BeginInstr = 1
	InsideInstr = 2
)

type InstrGroup struct {
	Label      string
	InstrList  []Instr
	PC         int
	LineNumber int
}

func PrintErrContext(tokensInLine []string, lineNum int, tokenNum int) {
	fmt.Printf("Error at Line %d, Token %d\n", lineNum, tokenNum)
	for _, token := range tokensInLine {
		fmt.Printf("%s ", token)
	}
	fmt.Println("")
	for i, token := range tokensInLine {
		if i==tokenNum {
			fmt.Printf("%s ", strings.Repeat("^", len(token)))
		} else {
			fmt.Printf("%s ", strings.Repeat(" ", len(token)))
		}
	}
	fmt.Println("")
}

func AssignPC(instrGroupList []InstrGroup) {
	pc := 0
	for _, instrGroup := range instrGroupList {
		instrGroup.PC = pc
		pc += len(instrGroup.InstrList)
	}
}

func BrTargetToImm(instrGroupList []InstrGroup) {
	targetPCMap := make(map[string]int)
	for _, instrGroup := range instrGroupList {
		if len(instrGroup.Label) != 0 {
			targetPCMap[instrGroup.Label] = instrGroup.PC
		}
	}
	for _, instrGroup := range instrGroupList {
		for _, instr := range instrGroup.InstrList {
			if len(instr.BrTarget) == 0 {
				continue
			}
			target, ok := targetPCMap[instr.BrTarget]
			if !ok {
				fmt.Printf("Line %d: Error! Cannot find such a lable: '%s'\n", instrGroup.LineNumber, instr.BrTarget)
				panic("Invalid Branch Target")
			}
			instr.Imm = int16(target - instrGroup.PC)
		}
	}
}

func ParseFile(fname string) (res []InstrGroup, ok bool) {
	tokensInLines := ReadLinesWithoutComments(fname)
	status := BeginGroup
	var currGroup InstrGroup
	var currInstr []string
	for lineNum, tokensInLine := range tokensInLines {
		for tokenNum, token := range tokensInLine {
			if status == BeginGroup {
				if token[:1] == "@" {
					currGroup.Label = token[1:]
				} else if token != "(G" {
					PrintErrContext(tokensInLine, lineNum, tokenNum)
					fmt.Printf("Expecting '(G' or '@', but meet '%s'\n", token)
					return
				} else {
					currGroup.LineNumber = lineNum
					status = BeginInstr
				}
			} else if status == BeginInstr {
				if token[:1] != "(" {
					PrintErrContext(tokensInLine, lineNum, tokenNum)
					fmt.Printf("Expecting '(' to start an instruction, but meet '%s'\n", token)
					return
				}
				if token=="(ret)" || token=="(clearbool)" || token=="(setbool)" {
					instr, _ := ParseInstr([]string{token[1:len(token)-1]}, lineNum)
					currGroup.InstrList = append(currGroup.InstrList, instr)
					status = BeginGroup
				} else {
					formatType := GetFormatType(token[1:])
					if formatType == Unknown {
						PrintErrContext(tokensInLine, lineNum, tokenNum)
						fmt.Printf("Not a valid mnemonic for instruction: '%s'\n", token[1:])
						return
					}
					currInstr = append(currInstr, token[1:])
					status = InsideInstr
				}
			} else if status == InsideInstr {
				endInstr := false
				endGroup := false
				endPos := 0
				if len(token) >= 2 && token[len(token)-2:] == "))" {
					endInstr = true
					endPos = len(token)-2
					endGroup = true
					status = BeginGroup
				} else if len(token) >= 1 && token[len(token)-1:] == ")" {
					endInstr = true
					endPos = len(token)-1
					status = BeginInstr
				} else {
					currInstr = append(currInstr, token)
				}
				if endInstr {
					if len(token[:endPos]) != 0 {
						currInstr = append(currInstr, token[:endPos])
					}
					instr, errMsg := ParseInstr(currInstr, lineNum)
					if len(errMsg) != 0 {
						PrintErrContext(tokensInLine, lineNum, tokenNum)
						fmt.Println(errMsg)
						return
					}
					currGroup.InstrList = append(currGroup.InstrList, instr)
					currInstr = currInstr[:0]
				}
				if endGroup {
					res = append(res, currGroup)
					currGroup.InstrList = nil
					currGroup.Label = ""
				}
			} else {
				panic("Cannot reach here!")
			}
		}
	}
	AssignPC(res)
	BrTargetToImm(res)
	ok = true
	return
}

func ReadLinesWithoutComments(fname string) (res [][]string) {
	file, err := os.Open(fname)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		poundPos := strings.Index(line, "#")
		if poundPos!=-1 {
			line = line[:poundPos]
		}
		line = strings.ToLower(line)
		line = strings.ReplaceAll(line, "\t", " ")
		tokens := strings.Split(line, " ")
		nonEmpty := make([]string, 0, len(tokens))
		for _, token := range tokens {
			if len(token) == 0 {
				continue
			}
			nonEmpty = append(nonEmpty, token)
		}
		res = append(res, nonEmpty)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return
}

func getRegID(token string, lineNum int) int8 {
	if len(token) <= 1 || (token[:1] != "r" && token[:1] != "R") {
		fmt.Printf("Line %d: Invalid RegID: %s\n", lineNum, token)
		panic("Invalid RegID")
	}
	id, err := strconv.Atoi(token[1:])
	if err != nil || id < 0 || id > 31 {
		fmt.Printf("Line %d: Invalid RegID: %s\n", lineNum, token)
		panic("Invalid RegID")
	}
	return int8(id)
}

func getImm(length int, token string, lineNum int) int16 {
	immTmp, err := strconv.Atoi(token)
	imm := int64(immTmp)
	if err!=nil {
		fmt.Printf("Line %d: Error when parsing '%s'\n", lineNum, token)
		panic(err)
	}
	usefulBits := 1
	for tmp := imm; tmp != 0 && tmp != -1; usefulBits++ {
		tmp = tmp>>1
	}
	if usefulBits > length {
		fmt.Printf("Line %d: '%s' cannot fit in %d bits\n", lineNum, token, length)
		panic(err)
	}
	return int16(imm)
}

func ParseInstr(instrTokens []string, lineNum int) (instr Instr, errMsg string) {
	instr.Mnemonic = instrTokens[0]
	switch GetFormatType(instr.Mnemonic) {
	case ControlFlow:
		if len(instrTokens)!=2 {
			errMsg = fmt.Sprintf("Invalid format! Should be: %s <label>", instr.Mnemonic)
			break
		}
		instr.BrTarget = instrTokens[1][1:]
	case NoParam:
		if len(instrTokens)!=1 {
			errMsg = fmt.Sprintf("Invalid format! Should have no parameters after %s", instr.Mnemonic)
			break
		}
	case SrcRegImm10:
		if len(instrTokens)!=3 {
			errMsg = fmt.Sprintf("Invalid format! Should be: %s <reg-id> <imm10>", instr.Mnemonic)
			break
		}
		instr.RD = getRegID(instrTokens[1], lineNum)
		instr.Imm = getImm(10, instrTokens[2], lineNum)
	case SrcRegImm8:
		if len(instrTokens)!=4 {
			errMsg = fmt.Sprintf("Invalid format! Should be: %s <dst-reg-id> <src-reg-id> <imm8>", instr.Mnemonic)
			break
		}
		instr.RD = getRegID(instrTokens[1], lineNum)
		instr.RB = getRegID(instrTokens[2], lineNum)
		instr.Imm = getImm(8, instrTokens[3], lineNum)
	case SrcRegX2:
		if len(instrTokens)!=4 {
			errMsg = fmt.Sprintf("Invalid format! Should be: %s <dst-reg-id> <src-reg0-id> <src-reg1-id>", instr.Mnemonic)
			break
		}
		instr.RD = getRegID(instrTokens[1], lineNum)
		instr.RA = getRegID(instrTokens[2], lineNum)
		instr.RB = getRegID(instrTokens[3], lineNum)
	case LoadRegIMM5:
		if len(instrTokens)!=4 {
			errMsg = fmt.Sprintf("Invalid format! Should be: %s <dst-reg-id> <addr-reg-id> <imm5>", instr.Mnemonic)
			break
		}
		instr.RD = getRegID(instrTokens[1], lineNum)
		instr.RB = getRegID(instrTokens[2], lineNum)
		instr.Imm = getImm(5, instrTokens[3], lineNum)
	case StoreRegIMM5:
		if len(instrTokens)!=4 {
			errMsg = fmt.Sprintf("Invalid format! Should be: %s <dst-reg-id> <addr-reg-id> <imm5>", instr.Mnemonic)
			break
		}
		instr.RD = getRegID(instrTokens[1], lineNum)
		instr.RB = getRegID(instrTokens[2], lineNum)
		instr.Imm = getImm(5, instrTokens[3], lineNum)
	case OnlyDst:
		if len(instrTokens)!=2 {
			errMsg = fmt.Sprintf("Invalid format! Should be: %s <dst-reg-id>", instr.Mnemonic)
			break
		}
		instr.RD = getRegID(instrTokens[1], lineNum)
	case SrcRegX1:
		if len(instrTokens)!=3 {
			errMsg = fmt.Sprintf("Invalid format! Should be: %s <dst-reg-id> <src-reg-id>", instr.Mnemonic)
			break
		}
		instr.RD = getRegID(instrTokens[1], lineNum)
		instr.RB = getRegID(instrTokens[2], lineNum)
	default:
		panic("Unknown FormatType")
	}
	return
}

func FormatInstr(instr Instr) string {
	var strb strings.Builder
	strb.WriteString("("+instr.Mnemonic)
	switch GetFormatType(instr.Mnemonic) {
	case ControlFlow:
		strb.WriteString(" "+instr.BrTarget+")")
	case NoParam:
		strb.WriteString(")")
	case SrcRegImm10:
		strb.WriteString(fmt.Sprintf(" r%d %d)", instr.RD, instr.Imm))
	case SrcRegImm8:
		strb.WriteString(fmt.Sprintf(" r%d r%d %d)", instr.RD, instr.RB, instr.Imm))
	case SrcRegX2:
		strb.WriteString(fmt.Sprintf(" r%d r%d r%d)", instr.RD, instr.RB, instr.RA))
	case LoadRegIMM5:
		strb.WriteString(fmt.Sprintf(" r%d r%d %d)", instr.RD, instr.RB, instr.Imm))
	case StoreRegIMM5:
		strb.WriteString(fmt.Sprintf(" r%d r%d %d)", instr.RD, instr.RB, instr.Imm))
	case OnlyDst:
		strb.WriteString(fmt.Sprintf(" r%d)", instr.RD))
	case SrcRegX1:
		strb.WriteString(fmt.Sprintf(" r%d r%d)", instr.RD, instr.RB))
	default:
		panic("Unknown FormatType")
	}
	return strb.String()
}

type FormatType int

const (
	Unknown FormatType = iota
	ControlFlow
	NoParam
	SrcRegImm10
	SrcRegImm8
	SrcRegX2
	LoadRegIMM5
	StoreRegIMM5
	OnlyDst
	SrcRegX1
)

func GetFormatType(mnemonic string) FormatType {
	switch mnemonic {
	case
	"br_ds_if0",
	"br_ds_if1",
	"br_nds_if0",
	"br_nds_if1",
	"jmp",
	"call":
		return ControlFlow
	case
	"ret",
	"clearbool",
	"setbool":
		return NoParam
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
		return SrcRegImm10
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
		return SrcRegImm8
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
	"mul_64",

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
	"xnmsbb_u",

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
	"f.sge_64",
	"f.madd_32",
	"f.madd_64",
	"f.msub_32",
	"f.msub_64",
	"f.nmadd_32",
	"f.nmadd_64",
	"f.nmsub_32",
	"f.nmsub_64":
		return SrcRegX2
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
		return LoadRegIMM5
	case
	"store128",
	"store_32",
	"store_64",
	"f.store_32",
	"f.store_64":
		return StoreRegIMM5
	case
	"yarnid_32",
	"yarnid_64":
		return OnlyDst
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
	"popcnt_64",
	"mov64loto32",
	"mov64hito32",
	"mov32to64_s",
	"mov32to64_u",

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
	"f.u32tof32",
	"f.u32tof64",
	"f.i32tof32",
	"f.i32tof64",
	"f.u64tof32",
	"f.u64tof64",
	"f.i64tof32",
	"f.i64tof64",
	"f.f32tou32",
	"f.f64tou32",
	"f.f32toi32",
	"f.f64toi32",
	"f.f32tou64",
	"f.f64tou64",
	"f.f32toi64",
	"f.f64toi64",
	"f.f32tof64",
	"f.f64tof32":
		return SrcRegX1
	default:
		return Unknown
	}
}

