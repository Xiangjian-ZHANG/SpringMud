// D C B A
//ra32Id
//useRa32
//rb32Id
//useRb32
//rc32Id
//useRc32
//rd32ID
//imm16
//ra64Id
//useRa64
//rb64Id
//useRb64
//rc64Id
//useRc64
//useRx64
//rd64ID
//
//useLL32
//useLL64
//
//simple32OP
//simple64OP
//complexOP
//memOP
//brOP

const bit[8:0]	BR_DS_IF0	= 9'h000;
const bit[8:0]	BR_DS_IF1	= 9'h001;
const bit[8:0]	BR_NDS_IF0	= 9'h002;
const bit[8:0]	BR_NDS_IF1	= 9'h003;
const bit[8:0]	JMP		= 9'h004;
const bit[8:0]	CALL		= 9'h005;
const bit[8:0]	AND_32		= 9'h006;
const bit[8:0]	AND_64		= 9'h007;
const bit[8:0]	OR_32		= 9'h008;
const bit[8:0]	OR_64		= 9'h009;
const bit[8:0]	XOR_32		= 9'h00A;
const bit[8:0]	XOR_64		= 9'h00B;
const bit[8:0]	UNISRC_32	= 9'h00C;
const bit[8:0]	UNISRC_64	= 9'h00D;
const bit[8:0]	SHL_32		= 9'h00E;
const bit[8:0]	SHL_64		= 9'h00F;
const bit[8:0]	SHR_S_32	= 9'h010;
const bit[8:0]	SHR_S_64	= 9'h011;
const bit[8:0]	SHR_U_32	= 9'h012;
const bit[8:0]	SHR_U_32	= 9'h013;
const bit[8:0]	ROTL_32		= 9'h014;
const bit[8:0]	ROTL_64		= 9'h015;
const bit[8:0]	ROTR_32		= 9'h016;
const bit[8:0]	ROTR_64		= 9'h017;
const bit[8:0]	SLT_S_32	= 9'h018;
const bit[8:0]	SLT_S_64	= 9'h019;
const bit[8:0]	SLT_U_32	= 9'h01A;
const bit[8:0]	SLT_U_64	= 9'h01B;
const bit[8:0]	SGE_S_32	= 9'h01C;
const bit[8:0]	SGE_S_64	= 9'h01D;
const bit[8:0]	SGE_U_32	= 9'h01E;
const bit[8:0]	SGE_U_64	= 9'h01F;
const bit[8:0]	SEQ_32		= 9'h020;
const bit[8:0]	SEQ_64		= 9'h021;
const bit[8:0]	SNE_32		= 9'h022;
const bit[8:0]	SNE_64		= 9'h023;
const bit[8:0]	BSLTI_S_32	= 9'h024;
const bit[8:0]	BSLTI_S_64	= 9'h025;
const bit[8:0]	BSLTI_U_32	= 9'h026;
const bit[8:0]	BSLTI_U_64	= 9'h027;
const bit[8:0]	BSGEI_S_32	= 9'h028;
const bit[8:0]	BSGEI_S_64	= 9'h029;
const bit[8:0]	BSGEI_U_32	= 9'h02A;
const bit[8:0]	BSGEI_U_64	= 9'h02B;
const bit[8:0]	BSEQI_32	= 9'h02C;
const bit[8:0]	BSEQI_64	= 9'h02D;
const bit[8:0]	BSNEI_32	= 9'h02E;
const bit[8:0]	BSNEI_64	= 9'h02F;
const bit[8:0]	ADDI_32		= 9'h030;
const bit[8:0]	ADDI_64		= 9'h031;
const bit[8:0]	SUBI_32		= 9'h040;
const bit[8:0]	SUBI_64		= 9'h041;
const bit[8:0]	ANDI_32		= 9'h050;
const bit[8:0]	ANDI_64		= 9'h051;
const bit[8:0]	ORI_32		= 9'h060;
const bit[8:0]	ORI_64		= 9'h061;
const bit[8:0]	XORI_32		= 9'h070;
const bit[8:0]	XORI_64		= 9'h071;
const bit[8:0]	SHLI_32		= 9'h080;
const bit[8:0]	SHLI_64		= 9'h081;
const bit[8:0]	SHRI_32		= 9'h090;
const bit[8:0]	SHRI_64		= 9'h091;
const bit[8:0]	ROTLI_32	= 9'h0A0;
const bit[8:0]	ROTLI_64	= 9'h0A1;
const bit[8:0]	SELECT_32	= 9'h0B0;
const bit[8:0]	SELECT_64	= 9'h0B1;
const bit[8:0]	CMOV_32		= 9'h0B2;
const bit[8:0]	CMOV_64		= 9'h0B3;
const bit[8:0]	MIN_S_32	= 9'h0B4;
const bit[8:0]	MIN_S_64	= 9'h0B5;
const bit[8:0]	MIN_U_32	= 9'h0B6;
const bit[8:0]	MIN_U_64	= 9'h0B7;
const bit[8:0]	MAX_S_32	= 9'h0B8;
const bit[8:0]	MAX_S_64	= 9'h0B9;
const bit[8:0]	MAX_U_32	= 9'h0BA;
const bit[8:0]	MAX_U_64	= 9'h0BB;
const bit[8:0]	ADD_32		= 9'h0C0;
const bit[8:0]	ADD_64		= 9'h0C1;
const bit[8:0]	SUB_32		= 9'h0C2;
const bit[8:0]	SUB_64		= 9'h0C3;
const bit[8:0]	ADC_32		= 9'h0C4;
const bit[8:0]	ADC_64		= 9'h0C5;
const bit[8:0]	SBB_32		= 9'h0C6;
const bit[8:0]	SBB_64		= 9'h0C7;
const bit[8:0]	MUL_32		= 9'h0C8;
const bit[8:0]	MUL_64		= 9'h0C9;
const bit[8:0]	LOAD_32		= 9'h0D0;
const bit[8:0]	LOAD_64		= 9'h0D1;
const bit[8:0]	FLOAD_32	= 9'h0D2;
const bit[8:0]	FLOAD_64	= 9'h0D3;
const bit[8:0]	STORE_32	= 9'h0D4;
const bit[8:0]	STORE_64	= 9'h0D5;
const bit[8:0]	FSTORE_32	= 9'h0D6;
const bit[8:0]	FSTORE_64	= 9'h0D7;
const bit[8:0]	LOADIMM_S_32	= 9'h0D8;
const bit[8:0]	LOADIMM_S_64	= 9'h0D9;
const bit[8:0]	LOADIMM_U_32	= 9'h0DA;
const bit[8:0]	LOADIMM_U_64	= 9'h0DB;
const bit[8:0]	LOAD8_S_32	= 9'h0DC;
const bit[8:0]	LOAD8_S_64	= 9'h0DD;
const bit[8:0]	LOAD8_U_32	= 9'h0DE;
const bit[8:0]	LOAD8_U_64	= 9'h0DF;
const bit[8:0]	LOAD16_S_32	= 9'h0E0;
const bit[8:0]	LOAD16_S_64	= 9'h0E1;
const bit[8:0]	LOAD16_U_32	= 9'h0E2;
const bit[8:0]	LOAD16_U_64	= 9'h0E3;
const bit[8:0]	STORE8_32	= 9'h0E4;
const bit[8:0]	STORE8_64	= 9'h0E5;
const bit[8:0]	STORE16_32	= 9'h0E6;
const bit[8:0]	STORE16_64	= 9'h0E7;
const bit[8:0]	STORE32_64	= 9'h0E8;
const bit[8:0]	LOAD32_S_64	= 9'h0E9;
const bit[8:0]	LOAD32_U_64	= 9'h0EA;
const bit[8:0]	LOAD128_64	= 9'h0EB;
const bit[8:0]	STORE128_64	= 9'h0EC;
const bit[8:0]	XMUL_S		= 9'h0F0;
const bit[8:0]	XMUL_U		= 9'h0F1;
const bit[8:0]	XMADD_S		= 9'h0F0;
const bit[8:0]	XMADD_U		= 9'h0F1;
const bit[8:0]	XMSUB_S		= 9'h0F2;
const bit[8:0]	XMSUB_U		= 9'h0F3;
const bit[8:0]	XNMADD_S	= 9'h0F4;
const bit[8:0]	XNMADD_U	= 9'h0F5;
const bit[8:0]	XNMSUB_S	= 9'h0F6;
const bit[8:0]	XNMSUB_U	= 9'h0F7;
const bit[8:0]	XMADC_S		= 9'h0F8;
const bit[8:0]	XMADC_U		= 9'h0F9;
const bit[8:0]	XMSBB_S		= 9'h0FA;
const bit[8:0]	XMSBB_U		= 9'h0FB;
const bit[8:0]	XNMADC_S	= 9'h0FC;
const bit[8:0]	XNMADC_U	= 9'h0FD;
const bit[8:0]	XNMSBB_S	= 9'h0FE;
const bit[8:0]	XNMSBB_U	= 9'h0FF;
const bit[8:0]	FADD_32		= 9'h180;
const bit[8:0]	FADD_64		= 9'h181;
const bit[8:0]	FSUB_32		= 9'h182;
const bit[8:0]	FSUB_64		= 9'h183;
const bit[8:0]	FMUL_32		= 9'h184;
const bit[8:0]	FMUL_64		= 9'h185;
const bit[8:0]	FMIN_32		= 9'h186;
const bit[8:0]	FMIN_64		= 9'h187;
const bit[8:0]	FMAX_32		= 9'h188;
const bit[8:0]	FMAX_64		= 9'h189;
const bit[8:0]	FCOPYSIGN_32	= 9'h18A;
const bit[8:0]	FCOPYSIGN_64	= 9'h18B;
const bit[8:0]	FUNISRC_32	= 9'h18C;
const bit[8:0]	FUNISRC_64	= 9'h18D;
const bit[8:0]	FSEQ_32		= 9'h190;
const bit[8:0]	FSEQ_64		= 9'h191;
const bit[8:0]	FSNE_32		= 9'h192;
const bit[8:0]	FSNE_64		= 9'h193;
const bit[8:0]	FSLT_32		= 9'h194;
const bit[8:0]	FSLT_64		= 9'h195;
const bit[8:0]	FSGE_32		= 9'h196;
const bit[8:0]	FSGE_64		= 9'h197;
const bit[8:0]	FCVT_32		= 9'h198;
const bit[8:0]	FCVT_64		= 9'h199;
const bit[8:0]	FMADD_32	= 9'h1A0;
const bit[8:0]	FMADD_64	= 9'h1A1;
const bit[8:0]	FMSUB_32	= 9'h1A2;
const bit[8:0]	FMSUB_64	= 9'h1A3;
const bit[8:0]	FNMADD_32	= 9'h1A8;
const bit[8:0]	FNMADD_64	= 9'h1A9;
const bit[8:0]	FNMSUB_32	= 9'h1AA;
const bit[8:0]	FNMSUB_64	= 9'h1AB;

const bit[3:0] YARNID		= 4'h0;
const bit[3:0] NOT		= 4'h1;
const bit[3:0] NEG		= 4'h2;
const bit[3:0] CLZ		= 4'h3;
const bit[3:0] CTZ		= 4'h4;
const bit[3:0] POPCNT		= 4'h5;
const bit[3:0] MOV64LOTO32 	= 4'h6;
const bit[3:0] MOV32TO64_S 	= 4'h6;
const bit[3:0] MOV64HITO32  	= 4'h7;
const bit[3:0] MOV32TO64_U 	= 4'h7;
const bit[3:0] CLEARBOOL	= 4'h8;
const bit[3:0] SETBOOL		= 4'h9;

const bit[3:0] FFLOOR  		= 4'h0;
const bit[3:0] FCEIL  		= 4'h0;
const bit[3:0] FTRUNC      	= 4'h1;
const bit[3:0] FROUND    	= 4'h2;
const bit[3:0] FABS        	= 4'h3;
const bit[3:0] FNEG        	= 4'h4;

const bit[3:0] FROMU32   	= 4'h0;
const bit[3:0] FROMU64       	= 4'h1;
const bit[3:0] FROMI32   	= 4'h2;
const bit[3:0] FROMI64       	= 4'h3;
const bit[3:0] TOU32         	= 4'h4;
const bit[3:0] TOU64         	= 4'h5;
const bit[3:0] TOI32         	= 4'h6;
const bit[3:0] TOI64         	= 4'h7;
const bit[3:0] TOANOTHER     	= 4'h8;

