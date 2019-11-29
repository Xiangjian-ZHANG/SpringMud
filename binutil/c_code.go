package binutil

const CLEAR_WR_EN = `
for (YI=0; YI<YC; YI++) {
	WriteRD32[YI]=WriteRD64[YI]=WriteRD64_X[YI]=WriteRD32_M[YI]=WriteRD64_M[YI]=WriteRD64_MX[YI]=false;
}
`
const DO_WRITE = `
for (YI=0; YI<YC; YI++) {
	if(WriteRD32[YI]) RF32[YI][RD32ID] = RD32Value[YI];
	if(WriteRD32_M[YI]) RF32[YI][RD32ID_M] = RD32Value_M[YI];
	if(WriteRD64[YI]) RF64[YI][RD64ID] = RD64Value[YI];
	if(WriteRD64_X[YI]) RF64[YI][RD64ID+1] = RD64Value_X[YI];
	if(WriteRD64_M[YI]) RF64[YI][RD64ID_M] = RD64Value_M[YI];
	if(WriteRD64_MX[YI]) RF64[YI][RD64ID_M+1] = RD64Value_MX[YI];
}
`

const C_FOR_SIM = `
const int YC = 16;
const int ST_ENTRY_COUNT = 32;

// Pipeline registers
int8_t RA32ID;
int8_t RB32ID;
int8_t RP32ID;
int8_t RM32ID;

int8_t RA64ID;
int8_t RB64ID;
int8_t RC64ID;
int8_t RM64ID;

int8_t RD32ID;
int8_t RD64ID;
int8_t RD32ID_M;
int8_t RD64ID_M;

int16_t Imm16For32;
int16_t Imm16For64;
int16_t Imm16ForMem;

bool    WriteRD32[YC];
bool    WriteRD64[YC];
bool    WriteRD64_X[YC];
int32_t RD32Value[YC];
int64_t RD64Value[YC];
int64_t RD64Value_X[YC];

bool    WriteRD32_M[YC];
bool    WriteRD64_M[YC];
bool    WriteRD64_MX[YC];
int32_t RD32Value_M[YC];
int64_t RD64Value_M[YC];
int64_t RD64Value_MX[YC];

bool    ShouldBranch;

// State regsisters
int64_t LastLoad64[YC];
int64_t RF64[YC][32];
int32_t RF32[YC][32];
bool    BoolReg[YC];
int32_t FwdBrTarget[YC];
int32_t SavedFwdBrTarget[YC][ST_ENTRY_COUNT];
void*   LinkReg[ST_ENTRY_COUNT];
int     ST_TOP;
int32_t PC;

// Helper variables
int YI; //Yarn index

#define BACK_BR_DS_IF0_PREPARE() \
	ShouldBranch = !BoolReg[0];

#define BACK_BR_DS_IF1_PREPARE() \
	ShouldBranch = BoolReg[0];

#define BACK_BR_DS_EXECUTE(label) \
	if(ShouldBranch) {ShouldBranch = false; goto label;}

#define BACK_BR_NDS_IF0(label) \
	if(!BoolReg[0]) goto label;

#define BACK_BR_NDS_IF1(label) \
	if(BoolReg[0]) goto label;

//-----

#define FWD_BR_NDS_IF0(label) \
	if(!BoolReg[YI]) FwdBrTarget[YI]=targetPC;

#define FWD_BR_NDS_IF1(label) \
	if(BoolReg[YI]) FwdBrTarget[YI]=targetPC;

#define JMP(label) goto label;

#define CALL(label, returnLabel) \
	for(YI=0; YI<YC; YI++) {\
	SavedFwdBrTarget[YI][ST_TOP]=FwdBrTarget[YI]; \
	FwdBrTarget[YI]=-1; \
	} \
	LinkReg[ST_TOP] = &&returnLabel; \
	ST_TOP++; \
	goto label;

#define RET() \
	ST_TOP--; \
	for(YI=0; YI<YC; YI++) {\
	FwdBrTarget[YI] = SavedFwdBrTarget[YI][ST_TOP]; \
	} \
	goto *(LinkReg[ST_TOP]); \

inline void AND_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RA32ID] & RF32[YI][RB32ID];
}
inline void AND_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RA64ID] & RF64[YI][RB64ID];
}
inline void OR_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RA32ID] | RF32[YI][RB32ID];
}
inline void OR_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RA64ID] | RF64[YI][RB64ID];
}
inline void XOR_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RA32ID] ^ RF32[YI][RB32ID];
}
inline void XOR_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RA64ID] ^ RF64[YI][RB64ID];
}
inline void UNISRC_32() {
	switch(RA32ID) {
	case YARNID:
		WriteRD32[YI] = true;
		RD32Value[YI] = YI;
		break;
	case NOT:
		WriteRD32[YI] = true;
		RD32Value[YI] = ~RF32[YI][RB32ID];
		break;
	case NEG
		WriteRD32[YI] = true;
		RD32Value[YI] = -int32_t(RF32[YI][RB32ID]);
		break;
	case CLZ:
		WriteRD32[YI] = true;
		RD32Value[YI] = _lz_cnt_u32(RF32[YI][RB32ID]);
		break;
	case CTZ:
		WriteRD32[YI] = true;
		RD32Value[YI] = _tz_cnt_u32(RF32[YI][RB32ID]);
		break;
	case POPCNT:
		WriteRD32[YI] = true;
		RD32Value[YI] = _popcnt32(RF32[YI][RB32ID]);
		break;
	case MOV32TO64_S:
		WriteRD64[YI] = true;
		RD64Value[YI] = int64_t(int32_t(RF32[YI][RB32ID]));
		break;
	case MOV32TO64_U:
		WriteRD64[YI] = true;
		RD64Value[YI] = uint64_t(RF32[YI][RB32ID]);
		break;
	case CLEARBOOL:
		BoolReg[YI] = false;
		break;
	case SETBOOL:
		BoolReg[YI] = true;
		break;
	default: printf("Invalid UNISRC_32\n")
		break;
	}
}
inline void UNISRC_64() {
	switch(RA64ID) {
	case YARNID:
		WriteRD64[YI] = true;
		RD64Value[YI] = YI;
		break;
	case NOT:
		WriteRD64[YI] = true;
		RD64Value[YI] = ~RF64[YI][RB64ID];
		break;
	case NEG
		WriteRD64[YI] = true;
		RD64Value[YI] = -int64_t(RF64[YI][RB64ID]);
		break;
	case CLZ:
		WriteRD64[YI] = true;
		RD64Value[YI] = _lz_cnt_u64(RF64[YI][RB64ID]);
		break;
	case CTZ:
		WriteRD64[YI] = true;
		RD64Value[YI] = _tz_cnt_u64(RF64[YI][RB64ID]);
		break;
	case POPCNT:
		WriteRD64[YI] = true;
		RD64Value[YI] = _popcnt64(RF64[YI][RB64ID]);
		break;
	case MOV64LOTO32:
		WriteRD64[YI] = true;
		RD32Value[YI] = uint32_t(RF64[YI][RB64ID]);
		break;
	case MOV64HITO32:
		WriteRD64[YI] = true;
		RD32Value[YI] = uint32_t(RF64[YI][RB64ID]>>32);
		break;
	default: printf("Invalid UNISRC_64\n")
		break;
	}
}
inline void SHL_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RA32ID] << RF32[YI][RB32ID];
}
inline void SHL_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RA64ID] << RF64[YI][RB64ID];
}
inline void SHR_S_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = int32_t(RF32[YI][RA32ID]) >> RF32[YI][RB32ID];
}
inline void SHR_S_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = int64_t(RF64[YI][RA64ID]) >> RF64[YI][RB64ID];
}
inline void SHR_U_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RA32ID] >> RF32[YI][RB32ID];
}
inline void SHR_U_32() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RA64ID] >> RF64[YI][RB64ID];
}
inline void ROTL_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = _rotl(RF32[YI][RA32ID], RF32[YI][RB32ID]);
}
inline void ROTL_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = _rotl64(RF64[YI][RA64ID], RF64[YI][RB64ID]);
}
inline void ROTR_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = _rotr(RF32[YI][RA32ID], RF32[YI][RB32ID]);
}
inline void ROTR_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = _rotr64(RF64[YI][RA64ID], RF64[YI][RB64ID]);
}
inline void SLT_S_32() {
	uint32_t res = (int32_t(RF32[YI][RA32ID]) < int32_t(RF32[YI][RB32ID])) ? 1 : 0;
	if (RD32ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD32[YI] = true;
		RD32Value = res;
	}
}
inline void SLT_S_64() {
	uint64_t res = (int64_t(RF64[YI][RA64ID]) < int64_t(RF64[YI][RB64ID])) ? 1 : 0;
	if (RD64ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD64[YI] = true;
		RD64Value = res;
	}
}
inline void SLT_U_32() {
	uint32_t res = (RF32[YI][RA32ID] < RF32[YI][RB32ID]) ? 1 : 0;
	if (RD32ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD32[YI] = true;
		RD32Value = res;
	}
}
inline void SLT_U_64() {
	uint64_t res = (RF64[YI][RA64ID] < RF64[YI][RB64ID]) ? 1 : 0;
	if (RD64ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD64[YI] = true;
		RD64Value = res;
	}
}
inline void SGE_S_32() {
	uint32_t res = (int32_t(RF32[YI][RA32ID]) >= int32_t(RF32[YI][RB32ID])) ? 1 : 0;
	if (RD32ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD32[YI] = true;
		RD32Value = res;
	}
}
inline void SGE_S_64() {
	uint64_t res = (int64_t(RF64[YI][RA64ID]) >= int64_t(RF64[YI][RB64ID])) ? 1 : 0;
	if (RD64ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD64[YI] = true;
		RD64Value = res;
	}
}
inline void SGE_U_32() {
	uint32_t res = (RF32[YI][RA32ID] >= RF32[YI][RB32ID]) ? 1 : 0;
	if (RD32ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD32[YI] = true;
		RD32Value = res;
	}
}
inline void SGE_U_64() {
	uint64_t res = (RF64[YI][RA64ID] >= RF64[YI][RB64ID]) ? 1 : 0;
	if (RD64ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD64[YI] = true;
		RD64Value = res;
	}
}
inline void SEQ_32() {
	uint32_t res = (RF32[YI][RA32ID] == RF32[YI][RB32ID]) ? 1 : 0;
	if (RD32ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD32[YI] = true;
		RD32Value = res;
	}
}
inline void SEQ_64() {
	uint64_t res = (RF64[YI][RA64ID] == RF64[YI][RB64ID]) ? 1 : 0;
	if (RD64ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD64[YI] = true;
		RD64Value = res;
	}
}
inline void SNE_32() {
	uint32_t res = (RF32[YI][RA32ID] != RF32[YI][RB32ID]) ? 1 : 0;
	if (RD32ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD32[YI] = true;
		RD32Value = res;
	}
}
inline void SNE_64() {
	uint64_t res = (RF64[YI][RA64ID] != RF64[YI][RB64ID]) ? 1 : 0;
	if (RD64ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD64[YI] = true;
		RD64Value = res;
	}
}
inline void BSLTI_S_32() {
	BoolReg = int32_t(RF32[YI][RB32ID]) < int32_t(Imm16For32);
}
inline void BSLTI_S_64() {
	BoolReg = int64_t(RF64[YI][RB64ID]) < int64_t(Imm16For64);
}
inline void BSLTI_U_32() {
	BoolReg = RF32[YI][RB32ID] < uint32_t(Imm16For32);
}
inline void BSLTI_U_64() {
	BoolReg = RF64[YI][RB64ID] < uint64_t(Imm16For64);
}
inline void BSGEI_S_32() {
	BoolReg = int32_t(RF32[YI][RB32ID]) >= int32_t(Imm16For32);
}
inline void BSGEI_S_64() {
	BoolReg = int64_t(RF64[YI][RB64ID]) >= int64_t(Imm16For64);
}
inline void BSGEI_U_32() {
	BoolReg = RF32[YI][RB32ID] >= uint32_t(Imm16For32);
}
inline void BSGEI_U_64() {
	BoolReg = RF64[YI][RB64ID] >= uint64_t(Imm16For64);
}
inline void BSEQI_S_32() {
	BoolReg = int32_t(RF32[YI][RB32ID]) == int32_t(Imm16For32);
}
inline void BSEQI_S_64() {
	BoolReg = int64_t(RF64[YI][RB64ID]) == int64_t(Imm16For64);
}
inline void BSNEI_U_32() {
	BoolReg = RF32[YI][RB32ID] != uint32_t(Imm16For32);
}
inline void BSNEI_U_64() {
	BoolReg = RF64[YI][RB64ID] != uint64_t(Imm16For64);
}
inline void ADDI_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RB32ID] + int32_t(Imm16For32);
}
inline void ADDI_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RB64ID] + int64_t(Imm16For64);
}
inline void SUBI_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RB32ID] - int32_t(Imm16For32);
}
inline void SUBI_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RB64ID] - int64_t(Imm16For64);
}
inline void ANDI_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RB32ID] & int32_t(Imm16For32);
}
inline void ANDI_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RB64ID] & int64_t(Imm16For64);
}
inline void ORI_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RB32ID] | int32_t(Imm16For32);
}
inline void ORI_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RB64ID] | int64_t(Imm16For64);
}
inline void XORI_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RB32ID] ^ int32_t(Imm16For32);
}
inline void XORI_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RB64ID] ^ int64_t(Imm16For64);
}
inline void SHLI_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RB32ID] + int32_t(Imm16For32);
}
inline void SHLI_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RB64ID] << int64_t(Imm16For64);
}
inline void SHRI_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RB32ID] << int32_t(Imm16For32);
}
inline void SHRI_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RB64ID] >> int64_t(Imm16For64);
}
inline void ROTLI_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = _rotl(RF32[YI][RB32ID], int32_t(Imm16For32));
}
inline void ROTLI_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = _rotl64(RF64[YI][RB64ID], int64_t(Imm16For64));
}
inline void SELECT_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = BoolReg[YI]? RF32[YI][RA32ID] : RF32[YI][RB32ID];
}
inline void SELECT_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = BoolReg[YI]? RF64[YI][RA64ID] : RF64[YI][RB64ID];
}
inline void CMOV_32() {
	WriteRD32[YI] = RF32[YI][RA32ID] != 0;
	RD32Value[YI] = RF32[YI][RB32ID];
}
inline void CMOV_64() {
	WriteRD64[YI] = RF64[YI][RA64ID] != 0;
	RD64Value[YI] = RF64[YI][RB64ID];
}
inline void MIN_S_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = MIN(int32_t(RF32[YI][RA32ID]), int32_t(RF32[YI][RB32ID]));
}
inline void MIN_S_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = MIN(int64_t(RF64[YI][RA64ID]), int64_t(RF64[YI][RB64ID]));
}
inline void MIN_U_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = MIN(RF32[YI][RA32ID], RF32[YI][RB32ID]);
}
inline void MIN_U_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = MIN(RF64[YI][RA64ID], RF64[YI][RB64ID]);
}
inline void MAX_S_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = MAX(int32_t(RF32[YI][RA32ID]), int32_t(RF32[YI][RB32ID]));
}
inline void MAX_S_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = MAX(int64_t(RF64[YI][RA64ID]), int64_t(RF64[YI][RB64ID]));
}
inline void MAX_U_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = MAX(RF32[YI][RA32ID], RF32[YI][RB32ID]);
}
inline void MAX_U_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = MAX(RF64[YI][RA64ID], RF64[YI][RB64ID]);
}
inline void ADD_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RA32ID] + RF32[YI][RB32ID];
}
inline void ADD_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RA64ID] + RF64[YI][RB64ID];
}
inline void SUB_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RA32ID] - RF32[YI][RB32ID];
}
inline void SUB_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RA64ID] - RF64[YI][RB64ID];
}
inline void ADC_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RA32ID] + RF32[YI][RB32ID] + (BoolReg[YI]?1:0);
	BoolReg[YI] = RF32[YI][RA32ID] >= RD32Value[YI];
}
inline void ADC_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RA64ID] + RF64[YI][RB64ID] + (BoolReg[YI]?1:0);
	BoolReg[YI] = RF64[YI][RA64ID] >= RD64Value[YI];
}
inline void SBB_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RA32ID] - RF32[YI][RB32ID] - (BoolReg[YI]?1:0);
	BoolReg[YI] = RF32[YI][RA32ID] <= RD32Value[YI];
}
inline void SBB_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RA64ID] - RF64[YI][RB64ID] - (BoolReg[YI]?1:0);
	BoolReg[YI] = RF64[YI][RA64ID] <= RD64Value[YI];
}
inline void MUL_S_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = int32_t(RF32[YI][RA32ID]) * int32_t(RF32[YI][RB32ID]);
}
inline void MUL_S_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = int64_t(RF64[YI][RA64ID]) * int64_t(RF64[YI][RB64ID]);
}
inline void MUL_U_32() {
	WriteRD32[YI] = true;
	RD32Value[YI] = RF32[YI][RA32ID] * RF32[YI][RB32ID];
}
inline void MUL_U_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = RF64[YI][RA64ID] * RF64[YI][RB64ID];
}
inline void LOAD_32() {
	WriteRD32_M[YI] = true;
	RD32Value_M[YI] = ((*uint32_t)Mem)[ RF32[YI][RP32ID] + int32_t(Imm16ForMem) ];
}
inline void LOAD_64() {
	WriteRD64_M[YI] = true;
	RD64Value_M[YI] = ((*uint64_t)Mem)[ RF32[YI][RP32ID] + int64_t(Imm16ForMem) ];
}
inline void FLOAD_32() {
	WriteRD64_M[YI] = true;
	RD64Value_M[YI] = ((*uint32_t)Mem)[ RF32[YI][RP32ID] + int32_t(Imm16ForMem) ];
}
inline void FLOAD_64() {
	WriteRD64_M[YI] = true;
	RD64Value_M[YI] = ((*uint64_t)Mem)[ RF32[YI][RP32ID] + int64_t(Imm16ForMem) ];
}
inline void STORE_32() {
	((*uint32_t)Mem)[ RF32[YI][RP32ID] + int32_t(Imm16ForMem) ] = RF32[YI][RM32ID];
}
inline void STORE_64() {
	((*uint64_t)Mem)[ RF32[YI][RP32ID] + int64_t(Imm16ForMem) ] = RF64[YI][RM64ID];
}
inline void FSTORE_32() {
	((*uint32_t)Mem)[ RF32[YI][RP32ID] + int32_t(Imm16ForMem) ] = (uint32_t)RF64[YI][RM32ID];
}
inline void FSTORE_64() {
	((*uint64_t)Mem)[ RF32[YI][RP32ID] + int64_t(Imm16ForMem) ] = RF64[YI][RM64ID];
}

inline void STORE8_32() {
	((*uint8_t)Mem)[ RF32[YI][RP32ID] + int32_t(Imm16ForMem) ] = RF32[YI][RM32ID];
}
inline void STORE8_64() {
	((*uint8_t)Mem)[ RF32[YI][RP32ID] + int64_t(Imm16ForMem) ] = RF64[YI][RM64ID];
}
inline void STORE16_32() {
	((*uint16_t)Mem)[ RF32[YI][RP32ID] + int32_t(Imm16ForMem) ] = RF32[YI][RM32ID];
}
inline void STORE16_64() {
	((*uint16_t)Mem)[ RF32[YI][RP32ID] + int64_t(Imm16ForMem) ] = RF64[YI][RM64ID];
}
inline void STORE32_64() {
	((*uint32_t)Mem)[ RF32[YI][RP32ID] + int64_t(Imm16ForMem) ] = RF64[YI][RM64ID];
}
inline void LOADIMM_S_32() {
	WriteRD32_M[YI] = true;
	RD32Value_M[YI] = int32_t(Imm16ForMem);
}
inline void LOADIMM_S_64() {
	WriteRD64_M[YI] = true;
	RD64Value_M[YI] = int64_t(Imm16ForMem);
}
inline void LOADIMM_U_32() {
	WriteRD32_M[YI] = true;
	RD32Value_M[YI] = uint32_t(Imm16ForMem);
}
inline void LOADIMM_U_64() {
	WriteRD64_M[YI] = true;
	RD64Value_M[YI] = uint64_t(Imm16ForMem);
}
inline void LOAD8_S_32() {
	WriteRD32_M[YI] = true;
	RD32Value_M[YI] = int32_t(((*int8_t)Mem)[ RF32[YI][RP32ID] + int32_t(Imm16ForMem) ]);
}
inline void LOAD8_S_64() {
	WriteRD64_M[YI] = true;
	RD64Value_M[YI] = int64_t(((*int8_t)Mem)[ RF32[YI][RP32ID] + int64_t(Imm16ForMem) ]);
}
inline void LOAD8_U_32() {
	WriteRD32_M[YI] = true;
	RD32Value_M[YI] = uint32_t(Mem[ RF32[YI][RP32ID] + int32_t(Imm16ForMem) ]);
}
inline void LOAD8_U_64() {
	WriteRD64_M[YI] = true;
	RD64Value_M[YI] = uint64_t(Mem[ RF32[YI][RP32ID] + int64_t(Imm16ForMem) ]);
}
inline void LOAD16_S_32() {
	WriteRD32_M[YI] = true;
	RD32Value_M[YI] = int32_t(((*int16_t)Mem)[ RF32[YI][RP32ID] + int32_t(Imm16ForMem) ]);
}
inline void LOAD16_S_64() {
	WriteRD64_M[YI] = true;
	RD64Value_M[YI] = int64_t(((*int16_t)Mem)[ RF32[YI][RP32ID] + int64_t(Imm16ForMem) ]);
}
inline void LOAD16_U_32() {
	WriteRD32_M[YI] = true;
	RD32Value_M[YI] = uint32_t(((*uint16_t)Mem)[ RF32[YI][RP32ID] + int32_t(Imm16ForMem) ]);
}
inline void LOAD16_U_64() {
	WriteRD64_M[YI] = true;
	RD64Value_M[YI] = uint64_t(((*uint16_t)Mem)[ RF32[YI][RP32ID] + int64_t(Imm16ForMem) ]);
}
inline void LOAD32_S_64() {
	WriteRD64_M[YI] = true;
	RD64Value_M[YI] = int64_t(((*int32_t)Mem)[ RF32[YI][RP32ID] + int64_t(Imm16ForMem) ]);
}
inline void LOAD32_U_64() {
	WriteRD64_M[YI] = true;
	RD64Value_M[YI] = uint64_t(((*uint32_t)Mem)[ RF32[YI][RP32ID] + int64_t(Imm16ForMem) ]);
}
inline void LOAD128_64() {
	WriteRD64_M[YI] = true;
	WriteRD64_MX[YI] = true;
	int addr = RF32[YI][RP32ID] + int64_t(Imm16ForMem);
	addr = (addr>>1)<<1; //clear lowest bit
	RD64Value_M[YI] = ((*uint64_t)Mem)[ addr ];
	RD64Value_MX[YI] = ((*uint64_t)Mem)[ addr+1 ];
}
inline void STORE128_64() {
	int addr = RF32[YI][RP32ID] + int64_t(Imm16ForMem);
	addr = (addr>>1)<<1; //clear lowest bit
	((*uint64_t)Mem)[ addr ] = RF64[YI][RM64ID];
	((*uint64_t)Mem)[ addr+1 ] = RF64[YI][RM64ID+1];
}
inline void xmul_func_s(bool tripleSrc, bool isAdd, bool needNeg, bool withCarry) {
	WriteRD64[YI] = true;
	WriteRD64_X[YI] = true;
	__int128_t tmp = 0;
	if (tripleSrc) {
		tmp = __int128_t(RF64[YI][RC64ID])<<64;
		tmp = tmp + RF64[YI][RC64ID+1];
	}
	__int128_t extra = (withCarry&&BoolReg[YI])? 1 :0;
	if (needNeg) {
		tmp = -tmp;
	}
	__int128_t op0 = __int128_t(int64_t(RF64[YI][RA64ID]));
	__int128_t op1 = __int128_t(int64_t(RF64[YI][RB64ID]));
	__int128_t res;
	bool carry;
	if (isAdd) {
		res = tmp + op0*op1 + extra;
		carry = res<tmp;
	} else {
		res = tmp - op0*op1 - extra;
		carry = res>tmp;
	}
	if (withCarry) {
		BoolReg[YI] = carry;
	}
	RD64Value[YI] = uint64_t(res);
	RD64Value_X[YI] = uint64_t(res>>64);
}
inline void xmul_func_u() {
	WriteRD64[YI] = true;
	WriteRD64_X[YI] = true;
	__uint128_t tmp = 0;
	if (tripleSrc) {
		tmp = __uint128_t(RF64[YI][RC64ID])<<64;
		tmp = tmp + RF64[YI][RC64ID+1];
	}
	__uint128_t extra = (withCarry&&BoolReg[YI])? 1 :0;
	if (needNeg) {
		tmp = -tmp;
	}
	__uint128_t op0 = __uint128_t(RF64[YI][RA64ID]);
	__uint128_t op1 = __uint128_t(RF64[YI][RB64ID]);
	__uint128_t res;
	bool carry;
	if (isAdd) {
		res = tmp + op0*op1 + extra;
		carry = res<tmp;
	} else {
		res = tmp - op0*op1 - extra;
		carry = res>tmp;
	}
	if (withCarry) {
		BoolReg[YI] = carry;
	}
	RD64Value[YI] = uint64_t(res);
	RD64Value_X[YI] = uint64_t(res>>64);
}
inline void XMUL_S() {
	xmul_func_s(false/*tripleSrc*/, false/*isAdd*/, false/*needNeg*/, false/*withCarry*/);
}
inline void XMUL_U() {
	xmul_func_u(false/*tripleSrc*/, false/*isAdd*/, false/*needNeg*/, false/*withCarry*/);
}
inline void XMADD_S() {
	xmul_func_s(true/*tripleSrc*/, true/*isAdd*/, false/*needNeg*/, false/*withCarry*/);
}
inline void XMADD_U() {
	xmul_func_u(true/*tripleSrc*/, true/*isAdd*/, false/*needNeg*/, false/*withCarry*/);
}
inline void XMSUB_S() {
	xmul_func_s(true/*tripleSrc*/, false/*isAdd*/, false/*needNeg*/, false/*withCarry*/);
}
inline void XMSUB_U() {
	xmul_func_u(true/*tripleSrc*/, false/*isAdd*/, false/*needNeg*/, false/*withCarry*/);
}
inline void XNMADD_S() {
	xmul_func_s(true/*tripleSrc*/, true/*isAdd*/, true/*needNeg*/, false/*withCarry*/);
}
inline void XNMADD_U() {
	xmul_func_u(true/*tripleSrc*/, true/*isAdd*/, true/*needNeg*/, false/*withCarry*/);
}
inline void XNMSUB_S() {
	xmul_func_s(true/*tripleSrc*/, false/*isAdd*/, true/*needNeg*/, false/*withCarry*/);
}
inline void XNMSUB_U() {
	xmul_func_u(true/*tripleSrc*/, false/*isAdd*/, true/*needNeg*/, false/*withCarry*/);
}
inline void XMADC_S() {
	xmul_func_s(true/*tripleSrc*/, true/*isAdd*/, false/*needNeg*/, true/*withCarry*/);
}
inline void XMADC_U() {
	xmul_func_u(true/*tripleSrc*/, true/*isAdd*/, false/*needNeg*/, true/*withCarry*/);
}
inline void XMSBB_S() {
	xmul_func_s(true/*tripleSrc*/, false/*isAdd*/, false/*needNeg*/, true/*withCarry*/);
}
inline void XMSBB_U() {
	xmul_func_u(true/*tripleSrc*/, false/*isAdd*/, false/*needNeg*/, true/*withCarry*/);
}
inline void XNMADC_S() {
	xmul_func_s(true/*tripleSrc*/, true/*isAdd*/, true/*needNeg*/, true/*withCarry*/);
}
inline void XNMADC_U() {
	xmul_func_u(true/*tripleSrc*/, true/*isAdd*/, true/*needNeg*/, true/*withCarry*/);
}
inline void XNMSBB_S() {
	xmul_func_s(true/*tripleSrc*/, false/*isAdd*/, true/*needNeg*/, true/*withCarry*/);
}
inline void XNMSBB_U() {
	xmul_func_u(true/*tripleSrc*/, false/*isAdd*/, true/*needNeg*/, true/*withCarry*/);
}
inline float FP32(uint32_t i) {
	uint32_t* ptr = &i;
	return *((float*)ptr);
}
inline uint32_t I32(float i) {
	*float ptr = &i;
	return *((uint32_t*)ptr);
}
inline double FP64(uint64_t i) {
	uint64_t* ptr = &i;
	return *((double*)ptr);
}
inline uint64_t I64(double i) {
	*double ptr = &i;
	return *((uint64_t*)ptr);
}
inline void FADD_32() {
	WriteRD64[YI] = true;
	RD64Value[YI] = FP32(RF64[YI][RA64ID]) + FP32(RF64[YI][RB64ID]);
}
inline void FADD_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = FP64(RF64[YI][RA64ID]) + FP64(RF64[YI][RB64ID]);
}
inline void FSUB_32() {
	WriteRD64[YI] = true;
	RD64Value[YI] = FP32(RF64[YI][RA64ID]) - FP32(RF64[YI][RB64ID]);
}
inline void FSUB_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = FP64(RF64[YI][RA64ID]) - FP64(RF64[YI][RB64ID]);
}
inline void FMUL_32() {
	WriteRD64[YI] = true;
	RD64Value[YI] = FP32(RF64[YI][RA64ID]) * FP32(RF64[YI][RB64ID]);
}
inline void FMUL_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = FP64(RF64[YI][RA64ID]) * FP64(RF64[YI][RB64ID]);
}
inline void FMIN_32() {
	WriteRD64[YI] = true;
	RD64Value[YI] = MIN(FP32(RF64[YI][RA64ID]), FP32(RF64[YI][RB64ID]));
}
inline void FMIN_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = MIN(FP64(RF64[YI][RA64ID]), FP64(RF64[YI][RB64ID]));
}
inline void FMAX_32() {
	WriteRD64[YI] = true;
	RD64Value[YI] = MAX(FP32(RF64[YI][RA64ID]), FP32(RF64[YI][RB64ID]));
}
inline void FMAX_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = MAX(FP64(RF64[YI][RA64ID]), FP64(RF64[YI][RB64ID]));
}
inline void FCOPYSIGN_32() {
	WriteRD64[YI] = true;
	uint32_t mask = 1<<31;
	RD64Value[YI] = (RF64[YI][RA64ID]&mask)|(RF64[YI][RB64ID]&(mask-1));
}
inline void FCOPYSIGN_64() {
	WriteRD64[YI] = true;
	uint64_t mask = uint64_t(1)<<63;
	RD64Value[YI] = (RF64[YI][RA64ID]&mask)|(RF64[YI][RB64ID]&~mask);
}
inline void FUNISRC_32() {
	switch(RA64ID) {
	case FFLOOR:
		WriteRD64[YI] = true;
		RD64Value[YI] = I32(floorf(FP32(RF64[YI][RA64ID])));
		break;
	case FTRUNC:
		WriteRD64[YI] = true;
		RD64Value[YI] = I32(truncf(FP32(RF64[YI][RA64ID])));
		break;
	case FROUND:
		WriteRD64[YI] = true;
		RD64Value[YI] = I32(roundf(FP32(RF64[YI][RA64ID])));
		break;
	case FABS:
		WriteRD64[YI] = true;
		RD64Value[YI] = I32(absf(FP32(RF64[YI][RA64ID])));
		break;
	case FNEG:
		WriteRD64[YI] = true;
		RD64Value[YI] = I32(-FP32(RF64[YI][RA64ID]));
		break;
	default: printf("Invalid FUNISRC_32\n")
		break;
	}
}
inline void FUNISRC_64() {
	switch(RA64ID) {
	case FFLOOR:
		WriteRD64[YI] = true;
		RD64Value[YI] = I64(floor(FP64(RF64[YI][RA64ID])));
		break;
	case FTRUNC:
		WriteRD64[YI] = true;
		RD64Value[YI] = I64(trunc(FP64(RF64[YI][RA64ID])));
		break;
	case FROUND:
		WriteRD64[YI] = true;
		RD64Value[YI] = I64(round(FP64(RF64[YI][RA64ID])));
		break;
	case FABS:
		WriteRD64[YI] = true;
		RD64Value[YI] = I64(abs(FP64(RF64[YI][RA64ID])));
		break;
	case FNEG:
		WriteRD64[YI] = true;
		RD64Value[YI] = I64(-FP64(RF64[YI][RA64ID]));
		break;
	default: printf("Invalid FUNISRC_64\n")
		break;
	}
}
inline void FSEQ_32() {
	uint64_t res = (FP32(RF64[YI][RA64ID]) == FP32(RF64[YI][RB64ID])) ? 1 : 0;
	if (RD64ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD64[YI] = true;
		RD64Value = res;
	}
}
inline void FSEQ_64() {
	uint64_t res = (FP64(RF64[YI][RA64ID]) == FP64(RF64[YI][RB64ID])) ? 1 : 0;
	if (RD64ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD64[YI] = true;
		RD64Value = res;
	}
}
inline void FSNE_32() {
	uint64_t res = (FP32(RF64[YI][RA64ID]) != FP32(RF64[YI][RB64ID])) ? 1 : 0;
	if (RD64ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD64[YI] = true;
		RD64Value = res;
	}
}
inline void FSNE_64() {
	uint64_t res = (FP64(RF64[YI][RA64ID]) != FP64(RF64[YI][RB64ID])) ? 1 : 0;
	if (RD64ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD64[YI] = true;
		RD64Value = res;
	}
}
inline void FSLT_32() {
	uint64_t res = (FP32(RF64[YI][RA64ID]) < FP32(RF64[YI][RB64ID])) ? 1 : 0;
	if (RD64ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD64[YI] = true;
		RD64Value = res;
	}
}
inline void FSLT_64() {
	uint64_t res = (FP64(RF64[YI][RA64ID]) < FP64(RF64[YI][RB64ID])) ? 1 : 0;
	if (RD64ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD64[YI] = true;
		RD64Value = res;
	}
}
inline void FSGE_32() {
	uint64_t res = (FP32(RF64[YI][RA64ID]) >= FP32(RF64[YI][RB64ID])) ? 1 : 0;
	if (RD64ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD64[YI] = true;
		RD64Value = res;
	}
}
inline void FSGE_64() {
	uint64_t res = (FP64(RF64[YI][RA64ID]) >= FP64(RF64[YI][RB64ID])) ? 1 : 0;
	if (RD64ID==0) {
		BoolReg = bool(res);
	} else {
		WriteRD64[YI] = true;
		RD64Value = res;
	}
}
inline void FCVT_32() {
	switch(RA64ID) {
	case FROMU32:
		WriteRD64[YI] = true;
		RD64Value[YI] = I32(float(RF32[YI][RA32ID]));
		break;
	case FROMU64:
		WriteRD64[YI] = true;
		RD64Value[YI] = I32(float(RF64[YI][RA64ID]));
		break;
	case FROMI32:
		WriteRD64[YI] = true;
		RD64Value[YI] = I32(float(int32_t(RF32[YI][RA32ID])));
		break;
	case FROMI64:
		WriteRD64[YI] = true;
		RD64Value[YI] = I32(float(int64_t(RF64[YI][RA64ID])));
		break;
	case TOU32:
		WriteRD32[YI] = true;
		RD32Value[YI] = uint32_t(FP32(RF64[YI][RA64ID]));
		break;
	case TOU64:
		WriteRD64[YI] = true;
		RD64Value[YI] = uint64_t(FP32(RF64[YI][RA64ID]));
		break;
	case TOI32:
		WriteRD32[YI] = true;
		RD32Value[YI] = int32_t(FP32(RF64[YI][RA64ID]));
		break;
	case TOI64:
		WriteRD64[YI] = true;
		RD64Value[YI] = int64_t(FP32(RF32[YI][RA32ID]));
		break;
	case TOANOTHER:
		WriteRD64[YI] = true;
		RD64Value[YI] = I64(double(FP32(RF64[YI][RA64ID])));
		break;
	default: printf("Invalid FCVT_32\n")
		break;
	}
}
inline void FCVT_64() {
	switch(RA32ID) {
	case FROMU32:
		WriteRD64[YI] = true;
		RD64Value[YI] = I64(double(RF32[YI][RA32ID]));
		break;
	case FROMU64:
		WriteRD64[YI] = true;
		RD64Value[YI] = I64(double(RF64[YI][RA64ID]));
		break;
	case FROMI32:
		WriteRD64[YI] = true;
		RD64Value[YI] = I64(double(int32_t(RF32[YI][RA32ID])));
		break;
	case FROMI64:
		WriteRD64[YI] = true;
		RD64Value[YI] = I64(double(int64_t(RF64[YI][RA64ID])));
		break;
	case TOU32:
		WriteRD32[YI] = true;
		RD32Value[YI] = int32_t(FP64(RF64[YI][RA64ID]));
		break;
	case TOU64:
		WriteRD64[YI] = true;
		RD64Value[YI] = uint64_t(FP64(RF64[YI][RA64ID]));
		break;
	case TOI32:
		WriteRD32[YI] = true;
		RD32Value[YI] = uint32_t(FP64(RF64[YI][RA64ID]));
		break;
	case TOI64:
		WriteRD64[YI] = true;
		RD64Value[YI] = int64_t(FP64(RF64[YI][RA64ID]));
		break;
	case TOANOTHER:
		WriteRD64[YI] = true;
		RD64Value[YI] = I32(float(FP64(RF64[YI][RA64ID])));
		break;
	default: printf("Invalid FCVT_64\n")
		break;
	}
}

//double fma  (double x     , double y     , double z);
//float fmaf (float x      , float y      , float z);
//Returns x*y+z.
inline void FMADD_32() {
	WriteRD64[YI] = true;
	RD64Value[YI] = I32(fma(FP32(RF64[YI][RA64ID]), FP32(RF64[YI][RB64ID]), FP32(RF64[YI][RC64ID])));
}
inline void FMADD_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = I64(fma(FP64(RF64[YI][RA64ID]), FP64(RF64[YI][RB64ID]), FP64(RF64[YI][RC64ID])));
}
inline void FMSUB_32() {
	WriteRD64[YI] = true;
	RD64Value[YI] = I32(fma(FP32(RF64[YI][RA64ID]), -FP32(RF64[YI][RB64ID]), FP32(RF64[YI][RC64ID])));
}
inline void FMSUB_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = I64(fma(FP64(RF64[YI][RA64ID]), -FP64(RF64[YI][RB64ID]), FP64(RF64[YI][RC64ID])));
}
inline void FNMADD_32() {
	WriteRD64[YI] = true;
	RD64Value[YI] = I32(fma(FP32(RF64[YI][RA64ID]), FP32(RF64[YI][RB64ID]), -FP32(RF64[YI][RC64ID])));
}
inline void FNMADD_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = I64(fma(FP64(RF64[YI][RA64ID]), FP64(RF64[YI][RB64ID]), -FP64(RF64[YI][RC64ID])));
}
inline void FNMSUB_32() {
	WriteRD64[YI] = true;
	RD64Value[YI] = I32(fma(FP32(RF64[YI][RA64ID]), -FP32(RF64[YI][RB64ID]), -FP32(RF64[YI][RC64ID])));
}
inline void FNMSUB_64() {
	WriteRD64[YI] = true;
	RD64Value[YI] = I64(fma(FP64(RF64[YI][RA64ID]), -FP64(RF64[YI][RB64ID]), -FP64(RF64[YI][RC64ID])));
}
`
