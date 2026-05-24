# OpenDART DS005 주요사항보고서 주요정보 — 양수도 Sub-1 «실물 양수도» 설계

- 작성일: 2026-05-24
- 모듈: `github.com/kenshin579/opendart`
- 범위: **DS005 양수도 그룹 중 Sub-1 8개 API** (`material` 패키지 확장)

## 배경 & 목표

DS005 부실·법적·증자·감자·사채 발행·자기주식(PR #10~#13)이 main 머지됨. DS005 양수도 그룹은
10개(영업/유형자산/타법인주식/주권사채권 양수·양도 8 + 자산양수도(기타)·풋백옵션 + 주식교환·이전 2)로
크기 때문에 **2개 sub-group 으로 분할**(A안). 이 spec 은 **Sub-1 «실물 양수도» 8개**(4쌍 양수/양도)다.
Sub-2(기타 2개)는 Sub-1 머지 후 별도 사이클. 기존 `MaterialParams` + `httpclient.GetList[T]` 재사용,
root 변경 없음. 신규 파일 `material/transfer.go`.

## API 표면 (docs 기반 사실)

- 8개 모두 동일 요청 `corp_code`+`bgn_de`+`end_de` (= `MaterialParams`), JSON `list[]`.
- 공통 블록: 머리(rcept_no/corp_cls/corp_code/corp_name) + 거래상대방(dlptn_*) + 외부평가(exevl_*) +
  거버넌스(bddd/od_a_at_t/b/adt_a_atn) + 공정위/풋옵션(ftc_stt_atn/popt_ctr_atn/popt_ctr_cn). 타입별 고유 필드는 각 struct.
- 값은 문자열(금액/주식수 콤마, 비율 %, 빈 값 "-").

## 아키텍처

```
material/
  transfer.go       # 8개 메서드 + item struct (신규)
  transfer_test.go  # 8개 fixture 테스트 (신규)
  testdata/         # 8개 fixture
README.md           # (수정) DS005 커버리지에 양수도(Sub-1)
integration_test.go # (수정) 통합 케이스 1~2개 (ErrNoData skip)
```

각 메서드: `func (c *Client) X(ctx, p MaterialParams) ([]XItem, error) { return httpclient.GetList[XItem](ctx, c.http, "<path>", p.toMap()) }`.

## 8개 메서드 (material/transfer.go)

| 메서드 | 한글 | 엔드포인트 | 필드 |
|--------|------|-----------|------|
| `BusinessAcquisition` | 영업양수 결정 | `/api/bsnInhDecsn.json` | 48 |
| `BusinessTransfer` | 영업양도 결정 | `/api/bsnTrfDecsn.json` | 41 |
| `TangibleAssetAcquisition` | 유형자산 양수 결정 | `/api/tgastInhDecsn.json` | 40 |
| `TangibleAssetTransfer` | 유형자산 양도 결정 | `/api/tgastTrfDecsn.json` | 40 |
| `OtherCorpStockAcquisition` | 타법인 주식·출자증권 양수결정 | `/api/otcprStkInvscrInhDecsn.json` | 43 |
| `OtherCorpStockTransfer` | 타법인 주식·출자증권 양도결정 | `/api/otcprStkInvscrTrfDecsn.json` | 39 |
| `StockRelatedBondAcquisition` | 주권 관련 사채권 양수 결정 | `/api/stkrtbdInhDecsn.json` | 41 |
| `StockRelatedBondTransfer` | 주권 관련 사채권 양도 결정 | `/api/stkrtbdTrfDecsn.json` | 41 |

```go
// BusinessAcquisitionItem 은 영업양수 결정 (bsnInhDecsn) 한 건.
type BusinessAcquisitionItem struct {
	RceptNo         string `json:"rcept_no"`            // 접수번호
	CorpCls         string `json:"corp_cls"`            // 법인구분 (Y/K/N/E)
	CorpCode        string `json:"corp_code"`           // 고유번호
	CorpName        string `json:"corp_name"`           // 회사명
	InhBsn          string `json:"inh_bsn"`             // 양수영업
	InhBsnMc        string `json:"inh_bsn_mc"`          // 양수영업 주요내용
	InhPrc          string `json:"inh_prc"`             // 양수가액(원)
	AbsnInhAtn      string `json:"absn_inh_atn"`        // 영업전부의 양수 여부
	AstInhBsn       string `json:"ast_inh_bsn"`         // 재무내용(자산액(양수대상 영업부문 A))
	AstCmpAll       string `json:"ast_cmp_all"`         // 재무내용(자산액(당사전체 B))
	AstRt           string `json:"ast_rt"`              // 재무내용(자산액 비중 % A/B)
	SlInhBsn        string `json:"sl_inh_bsn"`          // 재무내용(매출액(양수대상 영업부문 A))
	SlCmpAll        string `json:"sl_cmp_all"`          // 재무내용(매출액(당사전체 B))
	SlRt            string `json:"sl_rt"`               // 재무내용(매출액 비중 % A/B)
	DbtInhBsn       string `json:"dbt_inh_bsn"`         // 재무내용(부채액(양수대상 영업부문 A))
	DbtCmpAll       string `json:"dbt_cmp_all"`         // 재무내용(부채액(당사전체 B))
	DbtRt           string `json:"dbt_rt"`              // 재무내용(부채액 비중 % A/B)
	InhPp           string `json:"inh_pp"`              // 양수목적
	InhAf           string `json:"inh_af"`              // 양수영향
	InhPrdCtrCnsd   string `json:"inh_prd_ctr_cnsd"`    // 양수예정일자(계약체결일)
	InhPrdInhStd    string `json:"inh_prd_inh_std"`     // 양수예정일자(양수기준일)
	DlptnCmpnm      string `json:"dlptn_cmpnm"`         // 거래상대방(회사명/성명)
	DlptnCpt        string `json:"dlptn_cpt"`           // 거래상대방(자본금(원))
	DlptnMbsn       string `json:"dlptn_mbsn"`          // 거래상대방(주요사업)
	DlptnHoadd      string `json:"dlptn_hoadd"`         // 거래상대방(본점소재지(주소))
	DlptnRlCmpn     string `json:"dlptn_rl_cmpn"`       // 거래상대방(회사와의 관계)
	InhPym          string `json:"inh_pym"`             // 양수대금지급
	ExevlAtn        string `json:"exevl_atn"`           // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs       string `json:"exevl_bs_rs"`         // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn       string `json:"exevl_intn"`          // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd         string `json:"exevl_pd"`            // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp         string `json:"exevl_op"`            // 외부평가에 관한 사항(외부평가 의견)
	GmtsckSpdAtn    string `json:"gmtsck_spd_atn"`      // 주주총회 특별결의 여부
	GmtsckPrd       string `json:"gmtsck_prd"`          // 주주총회 예정일자
	AprskhPlnprc    string `json:"aprskh_plnprc"`       // 주식매수청구권(매수예정가격)
	AprskhPymPlpdMth string `json:"aprskh_pym_plpd_mth"` // 주식매수청구권(지급예정시기, 지급방법)
	AprskhLmt       string `json:"aprskh_lmt"`          // 주식매수청구권(제한 관련 내용)
	AprskhCtref     string `json:"aprskh_ctref"`        // 주식매수청구권(계약에 미치는 효력)
	Bddd            string `json:"bddd"`                // 이사회결의일(결정일)
	OdAAtT          string `json:"od_a_at_t"`           // 사외이사 참석여부(참석(명))
	OdAAtB          string `json:"od_a_at_b"`           // 사외이사 참석여부(불참(명))
	AdtAAtn         string `json:"adt_a_atn"`           // 감사(사외이사가 아닌 감사위원) 참석여부
	BdlstAtn        string `json:"bdlst_atn"`           // 우회상장 해당 여부
	N6mTpaiPlann    string `json:"n6m_tpai_plann"`      // 향후 6월이내 제3자배정 증자 등 계획
	OtcprBdlstSfAtn string `json:"otcpr_bdlst_sf_atn"`  // 타법인의 우회상장 요건 충족여부
	FtcSttAtn       string `json:"ftc_stt_atn"`         // 공정거래위원회 신고대상 여부
	PoptCtrAtn      string `json:"popt_ctr_atn"`        // 풋옵션 등 계약 체결여부
	PoptCtrCn       string `json:"popt_ctr_cn"`         // 계약내용
}

// BusinessTransferItem 은 영업양도 결정 (bsnTrfDecsn) 한 건.
type BusinessTransferItem struct {
	RceptNo         string `json:"rcept_no"`            // 접수번호
	CorpCls         string `json:"corp_cls"`            // 법인구분 (Y/K/N/E)
	CorpCode        string `json:"corp_code"`           // 고유번호
	CorpName        string `json:"corp_name"`           // 회사명
	TrfBsn          string `json:"trf_bsn"`             // 양도영업
	TrfBsnMc        string `json:"trf_bsn_mc"`          // 양도영업 주요내용
	TrfPrc          string `json:"trf_prc"`             // 양도가액(원)
	AstTrfBsn       string `json:"ast_trf_bsn"`         // 재무내용(자산액(양도대상 영업부문 A))
	AstCmpAll       string `json:"ast_cmp_all"`         // 재무내용(자산액(당사전체 B))
	AstRt           string `json:"ast_rt"`              // 재무내용(자산액 비중 % A/B)
	SlTrfBsn        string `json:"sl_trf_bsn"`          // 재무내용(매출액(양도대상 영업부문 A))
	SlCmpAll        string `json:"sl_cmp_all"`          // 재무내용(매출액(당사전체 B))
	SlRt            string `json:"sl_rt"`               // 재무내용(매출액 비중 % A/B)
	TrfPp           string `json:"trf_pp"`              // 양도목적
	TrfAf           string `json:"trf_af"`              // 양도영향
	TrfPrdCtrCnsd   string `json:"trf_prd_ctr_cnsd"`    // 양도예정일자(계약체결일)
	TrfPrdTrfStd    string `json:"trf_prd_trf_std"`     // 양도예정일자(양도기준일)
	DlptnCmpnm      string `json:"dlptn_cmpnm"`         // 거래상대방(회사명/성명)
	DlptnCpt        string `json:"dlptn_cpt"`           // 거래상대방(자본금(원))
	DlptnMbsn       string `json:"dlptn_mbsn"`          // 거래상대방(주요사업)
	DlptnHoadd      string `json:"dlptn_hoadd"`         // 거래상대방(본점소재지(주소))
	DlptnRlCmpn     string `json:"dlptn_rl_cmpn"`       // 거래상대방(회사와의 관계)
	TrfPym          string `json:"trf_pym"`             // 양도대금지급
	ExevlAtn        string `json:"exevl_atn"`           // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs       string `json:"exevl_bs_rs"`         // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn       string `json:"exevl_intn"`          // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd         string `json:"exevl_pd"`            // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp         string `json:"exevl_op"`            // 외부평가에 관한 사항(외부평가 의견)
	GmtsckSpdAtn    string `json:"gmtsck_spd_atn"`      // 주주총회 특별결의 여부
	GmtsckPrd       string `json:"gmtsck_prd"`          // 주주총회 예정일자
	AprskhPlnprc    string `json:"aprskh_plnprc"`       // 주식매수청구권(매수예정가격)
	AprskhPymPlpdMth string `json:"aprskh_pym_plpd_mth"` // 주식매수청구권(지급예정시기, 지급방법)
	AprskhLmt       string `json:"aprskh_lmt"`          // 주식매수청구권(제한 관련 내용)
	AprskhCtref     string `json:"aprskh_ctref"`        // 주식매수청구권(계약에 미치는 효력)
	Bddd            string `json:"bddd"`                // 이사회결의일(결정일)
	OdAAtT          string `json:"od_a_at_t"`           // 사외이사 참석여부(참석(명))
	OdAAtB          string `json:"od_a_at_b"`           // 사외이사 참석여부(불참(명))
	AdtAAtn         string `json:"adt_a_atn"`           // 감사(사외이사가 아닌 감사위원) 참석여부
	FtcSttAtn       string `json:"ftc_stt_atn"`         // 공정거래위원회 신고대상 여부
	PoptCtrAtn      string `json:"popt_ctr_atn"`        // 풋옵션 등 계약 체결여부
	PoptCtrCn       string `json:"popt_ctr_cn"`         // 계약내용
}

// TangibleAssetAcquisitionItem 은 유형자산 양수 결정 (tgastInhDecsn) 한 건.
type TangibleAssetAcquisitionItem struct {
	RceptNo           string `json:"rcept_no"`              // 접수번호
	CorpCls           string `json:"corp_cls"`              // 법인구분 (Y/K/N/E)
	CorpCode          string `json:"corp_code"`             // 고유번호
	CorpName          string `json:"corp_name"`             // 회사명
	AstSen            string `json:"ast_sen"`               // 자산구분
	AstNm             string `json:"ast_nm"`                // 자산명
	InhdtlInhprc      string `json:"inhdtl_inhprc"`         // 양수내역(양수금액(원))
	InhdtlTast        string `json:"inhdtl_tast"`           // 양수내역(자산총액(원))
	InhdtlTastVs      string `json:"inhdtl_tast_vs"`        // 양수내역(자산총액대비(%))
	InhPp             string `json:"inh_pp"`                // 양수목적
	InhAf             string `json:"inh_af"`                // 양수영향
	InhPrdCtrCnsd     string `json:"inh_prd_ctr_cnsd"`      // 양수예정일자(계약체결일)
	InhPrdInhStd      string `json:"inh_prd_inh_std"`       // 양수예정일자(양수기준일)
	InhPrdRgsPrd      string `json:"inh_prd_rgs_prd"`       // 양수예정일자(등기예정일)
	DlptnCmpnm        string `json:"dlptn_cmpnm"`           // 거래상대방(회사명/성명)
	DlptnCpt          string `json:"dlptn_cpt"`             // 거래상대방(자본금(원))
	DlptnMbsn         string `json:"dlptn_mbsn"`            // 거래상대방(주요사업)
	DlptnHoadd        string `json:"dlptn_hoadd"`           // 거래상대방(본점소재지(주소))
	DlptnRlCmpn       string `json:"dlptn_rl_cmpn"`         // 거래상대방(회사와의 관계)
	DlPym             string `json:"dl_pym"`                // 거래대금지급
	ExevlAtn          string `json:"exevl_atn"`             // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs         string `json:"exevl_bs_rs"`           // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn         string `json:"exevl_intn"`            // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd           string `json:"exevl_pd"`              // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp           string `json:"exevl_op"`              // 외부평가에 관한 사항(외부평가 의견)
	GmtsckSpdAtn      string `json:"gmtsck_spd_atn"`        // 주주총회 특별결의 여부
	GmtsckPrd         string `json:"gmtsck_prd"`            // 주주총회 예정일자
	AprskhExrq        string `json:"aprskh_exrq"`           // 주식매수청구권(행사요건)
	AprskhPlnprc      string `json:"aprskh_plnprc"`         // 주식매수청구권(매수예정가격)
	AprskhExPcMthPdPl string `json:"aprskh_ex_pc_mth_pd_pl"` // 주식매수청구권(행사절차, 방법, 기간, 장소)
	AprskhPymPlpdMth  string `json:"aprskh_pym_plpd_mth"`   // 주식매수청구권(지급예정시기, 지급방법)
	AprskhLmt         string `json:"aprskh_lmt"`            // 주식매수청구권(제한 관련 내용)
	AprskhCtref       string `json:"aprskh_ctref"`          // 주식매수청구권(계약에 미치는 효력)
	Bddd              string `json:"bddd"`                  // 이사회결의일(결정일)
	OdAAtT            string `json:"od_a_at_t"`             // 사외이사 참석여부(참석(명))
	OdAAtB            string `json:"od_a_at_b"`             // 사외이사 참석여부(불참(명))
	AdtAAtn           string `json:"adt_a_atn"`             // 감사(사외이사가 아닌 감사위원) 참석여부
	FtcSttAtn         string `json:"ftc_stt_atn"`           // 공정거래위원회 신고대상 여부
	PoptCtrAtn        string `json:"popt_ctr_atn"`          // 풋옵션 등 계약 체결여부
	PoptCtrCn         string `json:"popt_ctr_cn"`           // 계약내용
}

// TangibleAssetTransferItem 은 유형자산 양도 결정 (tgastTrfDecsn) 한 건.
type TangibleAssetTransferItem struct {
	RceptNo           string `json:"rcept_no"`              // 접수번호
	CorpCls           string `json:"corp_cls"`              // 법인구분 (Y/K/N/E)
	CorpCode          string `json:"corp_code"`             // 고유번호
	CorpName          string `json:"corp_name"`             // 회사명
	AstSen            string `json:"ast_sen"`               // 자산구분
	AstNm             string `json:"ast_nm"`                // 자산명
	TrfdtlTrfprc      string `json:"trfdtl_trfprc"`         // 양도내역(양도금액(원))
	TrfdtlTast        string `json:"trfdtl_tast"`           // 양도내역(자산총액(원))
	TrfdtlTastVs      string `json:"trfdtl_tast_vs"`        // 양도내역(자산총액대비(%))
	TrfPp             string `json:"trf_pp"`                // 양도목적
	TrfAf             string `json:"trf_af"`                // 양도영향
	TrfPrdCtrCnsd     string `json:"trf_prd_ctr_cnsd"`      // 양도예정일자(계약체결일)
	TrfPrdTrfStd      string `json:"trf_prd_trf_std"`       // 양도예정일자(양도기준일)
	TrfPrdRgsPrd      string `json:"trf_prd_rgs_prd"`       // 양도예정일자(등기예정일)
	DlptnCmpnm        string `json:"dlptn_cmpnm"`           // 거래상대방(회사명/성명)
	DlptnCpt          string `json:"dlptn_cpt"`             // 거래상대방(자본금(원))
	DlptnMbsn         string `json:"dlptn_mbsn"`            // 거래상대방(주요사업)
	DlptnHoadd        string `json:"dlptn_hoadd"`           // 거래상대방(본점소재지(주소))
	DlptnRlCmpn       string `json:"dlptn_rl_cmpn"`         // 거래상대방(회사와의 관계)
	DlPym             string `json:"dl_pym"`                // 거래대금지급
	ExevlAtn          string `json:"exevl_atn"`             // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs         string `json:"exevl_bs_rs"`           // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn         string `json:"exevl_intn"`            // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd           string `json:"exevl_pd"`              // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp           string `json:"exevl_op"`              // 외부평가에 관한 사항(외부평가 의견)
	GmtsckSpdAtn      string `json:"gmtsck_spd_atn"`        // 주주총회 특별결의 여부
	GmtsckPrd         string `json:"gmtsck_prd"`            // 주주총회 예정일자
	AprskhExrq        string `json:"aprskh_exrq"`           // 주식매수청구권(행사요건)
	AprskhPlnprc      string `json:"aprskh_plnprc"`         // 주식매수청구권(매수예정가격)
	AprskhExPcMthPdPl string `json:"aprskh_ex_pc_mth_pd_pl"` // 주식매수청구권(행사절차, 방법, 기간, 장소)
	AprskhPymPlpdMth  string `json:"aprskh_pym_plpd_mth"`   // 주식매수청구권(지급예정시기, 지급방법)
	AprskhLmt         string `json:"aprskh_lmt"`            // 주식매수청구권(제한 관련 내용)
	AprskhCtref       string `json:"aprskh_ctref"`          // 주식매수청구권(계약에 미치는 효력)
	Bddd              string `json:"bddd"`                  // 이사회결의일(결정일)
	OdAAtT            string `json:"od_a_at_t"`             // 사외이사 참석여부(참석(명))
	OdAAtB            string `json:"od_a_at_b"`             // 사외이사 참석여부(불참(명))
	AdtAAtn           string `json:"adt_a_atn"`             // 감사(사외이사가 아닌 감사위원) 참석여부
	FtcSttAtn         string `json:"ftc_stt_atn"`           // 공정거래위원회 신고대상 여부
	PoptCtrAtn        string `json:"popt_ctr_atn"`          // 풋옵션 등 계약 체결여부
	PoptCtrCn         string `json:"popt_ctr_cn"`           // 계약내용
}

// OtherCorpStockAcquisitionItem 은 타법인 주식 및 출자증권 양수결정 (otcprStkInvscrInhDecsn) 한 건.
type OtherCorpStockAcquisitionItem struct {
	RceptNo         string `json:"rcept_no"`            // 접수번호
	CorpCls         string `json:"corp_cls"`            // 법인구분 (Y/K/N/E)
	CorpCode        string `json:"corp_code"`           // 고유번호
	CorpName        string `json:"corp_name"`           // 회사명
	IscmpCmpnm      string `json:"iscmp_cmpnm"`         // 발행회사(회사명)
	IscmpNt         string `json:"iscmp_nt"`            // 발행회사(국적)
	IscmpRp         string `json:"iscmp_rp"`            // 발행회사(대표자)
	IscmpCpt        string `json:"iscmp_cpt"`           // 발행회사(자본금(원))
	IscmpRlCmpn     string `json:"iscmp_rl_cmpn"`       // 발행회사(회사와 관계)
	IscmpTisstk     string `json:"iscmp_tisstk"`        // 발행회사(발행주식 총수(주))
	IscmpMbsn       string `json:"iscmp_mbsn"`          // 발행회사(주요사업)
	L6mTpaNstkaqAtn string `json:"l6m_tpa_nstkaq_atn"`  // 최근 6월 이내 제3자 배정에 의한 신주취득 여부
	InhdtlStkcnt    string `json:"inhdtl_stkcnt"`       // 양수내역(양수주식수(주))
	InhdtlInhprc    string `json:"inhdtl_inhprc"`       // 양수내역(양수금액(원)(A))
	InhdtlTast      string `json:"inhdtl_tast"`         // 양수내역(총자산(원)(B))
	InhdtlTastVs    string `json:"inhdtl_tast_vs"`      // 양수내역(총자산대비(%)(A/B))
	InhdtlEcpt      string `json:"inhdtl_ecpt"`         // 양수내역(자기자본(원)(C))
	InhdtlEcptVs    string `json:"inhdtl_ecpt_vs"`      // 양수내역(자기자본대비(%)(A/C))
	AtinhOwstkcnt   string `json:"atinh_owstkcnt"`      // 양수후 소유주식수 및 지분비율(소유주식수(주))
	AtinhEqrt       string `json:"atinh_eqrt"`          // 양수후 소유주식수 및 지분비율(지분비율(%))
	InhPp           string `json:"inh_pp"`              // 양수목적
	InhPrd          string `json:"inh_prd"`             // 양수예정일자
	DlptnCmpnm      string `json:"dlptn_cmpnm"`         // 거래상대방(회사명/성명)
	DlptnCpt        string `json:"dlptn_cpt"`           // 거래상대방(자본금(원))
	DlptnMbsn       string `json:"dlptn_mbsn"`          // 거래상대방(주요사업)
	DlptnHoadd      string `json:"dlptn_hoadd"`         // 거래상대방(본점소재지(주소))
	DlptnRlCmpn     string `json:"dlptn_rl_cmpn"`       // 거래상대방(회사와의 관계)
	DlPym           string `json:"dl_pym"`              // 거래대금지급
	ExevlAtn        string `json:"exevl_atn"`           // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs       string `json:"exevl_bs_rs"`         // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn       string `json:"exevl_intn"`          // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd         string `json:"exevl_pd"`            // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp         string `json:"exevl_op"`            // 외부평가에 관한 사항(외부평가 의견)
	Bddd            string `json:"bddd"`                // 이사회결의일(결정일)
	OdAAtT          string `json:"od_a_at_t"`           // 사외이사 참석여부(참석(명))
	OdAAtB          string `json:"od_a_at_b"`           // 사외이사 참석여부(불참(명))
	AdtAAtn         string `json:"adt_a_atn"`           // 감사(사외이사가 아닌 감사위원) 참석여부
	BdlstAtn        string `json:"bdlst_atn"`           // 우회상장 해당 여부
	N6mTpaiPlann    string `json:"n6m_tpai_plann"`      // 향후 6월이내 제3자배정 증자 등 계획
	IscmpBdlstSfAtn string `json:"iscmp_bdlst_sf_atn"`  // 발행회사(타법인)의 우회상장 요건 충족여부
	FtcSttAtn       string `json:"ftc_stt_atn"`         // 공정거래위원회 신고대상 여부
	PoptCtrAtn      string `json:"popt_ctr_atn"`        // 풋옵션 등 계약 체결여부
	PoptCtrCn       string `json:"popt_ctr_cn"`         // 계약내용
}

// OtherCorpStockTransferItem 은 타법인 주식 및 출자증권 양도결정 (otcprStkInvscrTrfDecsn) 한 건.
type OtherCorpStockTransferItem struct {
	RceptNo       string `json:"rcept_no"`        // 접수번호
	CorpCls       string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`       // 고유번호
	CorpName      string `json:"corp_name"`       // 회사명
	IscmpCmpnm    string `json:"iscmp_cmpnm"`     // 발행회사(회사명)
	IscmpNt       string `json:"iscmp_nt"`        // 발행회사(국적)
	IscmpRp       string `json:"iscmp_rp"`        // 발행회사(대표자)
	IscmpCpt      string `json:"iscmp_cpt"`       // 발행회사(자본금(원))
	IscmpRlCmpn   string `json:"iscmp_rl_cmpn"`   // 발행회사(회사와 관계)
	IscmpTisstk   string `json:"iscmp_tisstk"`    // 발행회사(발행주식 총수(주))
	IscmpMbsn     string `json:"iscmp_mbsn"`      // 발행회사(주요사업)
	TrfdtlStkcnt  string `json:"trfdtl_stkcnt"`   // 양도내역(양도주식수(주))
	TrfdtlTrfprc  string `json:"trfdtl_trfprc"`   // 양도내역(양도금액(원)(A))
	TrfdtlTast    string `json:"trfdtl_tast"`     // 양도내역(총자산(원)(B))
	TrfdtlTastVs  string `json:"trfdtl_tast_vs"`  // 양도내역(총자산대비(%)(A/B))
	TrfdtlEcpt    string `json:"trfdtl_ecpt"`     // 양도내역(자기자본(원)(C))
	TrfdtlEcptVs  string `json:"trfdtl_ecpt_vs"`  // 양도내역(자기자본대비(%)(A/C))
	AttrfOwstkcnt string `json:"attrf_owstkcnt"`  // 양도후 소유주식수 및 지분비율(소유주식수(주))
	AttrfEqrt     string `json:"attrf_eqrt"`      // 양도후 소유주식수 및 지분비율(지분비율(%))
	TrfPp         string `json:"trf_pp"`          // 양도목적
	TrfPrd        string `json:"trf_prd"`         // 양도예정일자
	DlptnCmpnm    string `json:"dlptn_cmpnm"`     // 거래상대방(회사명/성명)
	DlptnCpt      string `json:"dlptn_cpt"`       // 거래상대방(자본금(원))
	DlptnMbsn     string `json:"dlptn_mbsn"`      // 거래상대방(주요사업)
	DlptnHoadd    string `json:"dlptn_hoadd"`     // 거래상대방(본점소재지(주소))
	DlptnRlCmpn   string `json:"dlptn_rl_cmpn"`   // 거래상대방(회사와의 관계)
	DlPym         string `json:"dl_pym"`          // 거래대금지급
	ExevlAtn      string `json:"exevl_atn"`       // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs     string `json:"exevl_bs_rs"`     // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn     string `json:"exevl_intn"`      // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd       string `json:"exevl_pd"`        // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp       string `json:"exevl_op"`        // 외부평가에 관한 사항(외부평가 의견)
	Bddd          string `json:"bddd"`            // 이사회결의일(결정일)
	OdAAtT        string `json:"od_a_at_t"`       // 사외이사 참석여부(참석(명))
	OdAAtB        string `json:"od_a_at_b"`       // 사외이사 참석여부(불참(명))
	AdtAAtn       string `json:"adt_a_atn"`       // 감사(사외이사가 아닌 감사위원) 참석여부
	FtcSttAtn     string `json:"ftc_stt_atn"`     // 공정거래위원회 신고대상 여부
	PoptCtrAtn    string `json:"popt_ctr_atn"`    // 풋옵션 등 계약 체결여부
	PoptCtrCn     string `json:"popt_ctr_cn"`     // 계약내용
}

// StockRelatedBondAcquisitionItem 은 주권 관련 사채권 양수 결정 (stkrtbdInhDecsn) 한 건.
type StockRelatedBondAcquisitionItem struct {
	RceptNo         string `json:"rcept_no"`            // 접수번호
	CorpCls         string `json:"corp_cls"`            // 법인구분 (Y/K/N/E)
	CorpCode        string `json:"corp_code"`           // 고유번호
	CorpName        string `json:"corp_name"`           // 회사명
	StkrtbdKndn     string `json:"stkrtbd_kndn"`        // 주권 관련 사채권의 종류
	Tm              string `json:"tm"`                  // 주권 관련 사채권의 종류(회차)
	Knd             string `json:"knd"`                 // 주권 관련 사채권의 종류(종류)
	BdiscmpCmpnm    string `json:"bdiscmp_cmpnm"`       // 사채권 발행회사(회사명)
	BdiscmpNt       string `json:"bdiscmp_nt"`          // 사채권 발행회사(국적)
	BdiscmpRp       string `json:"bdiscmp_rp"`          // 사채권 발행회사(대표자)
	BdiscmpCpt      string `json:"bdiscmp_cpt"`         // 사채권 발행회사(자본금(원))
	BdiscmpRlCmpn   string `json:"bdiscmp_rl_cmpn"`     // 사채권 발행회사(회사와 관계)
	BdiscmpTisstk   string `json:"bdiscmp_tisstk"`      // 사채권 발행회사(발행주식 총수(주))
	BdiscmpMbsn     string `json:"bdiscmp_mbsn"`        // 사채권 발행회사(주요사업)
	L6mTpaNstkaqAtn string `json:"l6m_tpa_nstkaq_atn"`  // 최근 6월 이내 제3자 배정에 의한 신주취득 여부
	InhdtlBdFta     string `json:"inhdtl_bd_fta"`       // 양수내역(사채의 권면(전자등록)총액(원))
	InhdtlInhprc    string `json:"inhdtl_inhprc"`       // 양수내역(양수금액(원)(A))
	InhdtlTast      string `json:"inhdtl_tast"`         // 양수내역(총자산(원)(B))
	InhdtlTastVs    string `json:"inhdtl_tast_vs"`      // 양수내역(총자산대비(%)(A/B))
	InhdtlEcpt      string `json:"inhdtl_ecpt"`         // 양수내역(자기자본(원)(C))
	InhdtlEcptVs    string `json:"inhdtl_ecpt_vs"`      // 양수내역(자기자본대비(%)(A/C))
	InhPp           string `json:"inh_pp"`              // 양수목적
	InhPrd          string `json:"inh_prd"`             // 양수예정일자
	DlptnCmpnm      string `json:"dlptn_cmpnm"`         // 거래상대방(회사명/성명)
	DlptnCpt        string `json:"dlptn_cpt"`           // 거래상대방(자본금(원))
	DlptnMbsn       string `json:"dlptn_mbsn"`          // 거래상대방(주요사업)
	DlptnHoadd      string `json:"dlptn_hoadd"`         // 거래상대방(본점소재지(주소))
	DlptnRlCmpn     string `json:"dlptn_rl_cmpn"`       // 거래상대방(회사와의 관계)
	DlPym           string `json:"dl_pym"`              // 거래대금지급
	ExevlAtn        string `json:"exevl_atn"`           // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs       string `json:"exevl_bs_rs"`         // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn       string `json:"exevl_intn"`          // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd         string `json:"exevl_pd"`            // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp         string `json:"exevl_op"`            // 외부평가에 관한 사항(외부평가 의견)
	Bddd            string `json:"bddd"`                // 이사회결의일(결정일)
	OdAAtT          string `json:"od_a_at_t"`           // 사외이사 참석여부(참석(명))
	OdAAtB          string `json:"od_a_at_b"`           // 사외이사 참석여부(불참(명))
	AdtAAtn         string `json:"adt_a_atn"`           // 감사(사외이사가 아닌 감사위원) 참석여부
	FtcSttAtn       string `json:"ftc_stt_atn"`         // 공정거래위원회 신고대상 여부
	PoptCtrAtn      string `json:"popt_ctr_atn"`        // 풋옵션 등 계약 체결여부
	PoptCtrCn       string `json:"popt_ctr_cn"`         // 계약내용
}

// StockRelatedBondTransferItem 은 주권 관련 사채권 양도 결정 (stkrtbdTrfDecsn) 한 건.
type StockRelatedBondTransferItem struct {
	RceptNo       string `json:"rcept_no"`        // 접수번호
	CorpCls       string `json:"corp_cls"`        // 법인구분 (Y/K/N/E)
	CorpCode      string `json:"corp_code"`       // 고유번호
	CorpName      string `json:"corp_name"`       // 회사명
	StkrtbdKndn   string `json:"stkrtbd_kndn"`    // 주권 관련 사채권의 종류
	Tm            string `json:"tm"`              // 주권 관련 사채권의 종류(회차)
	Knd           string `json:"knd"`             // 주권 관련 사채권의 종류(종류)
	Aqd           string `json:"aqd"`             // 취득일자
	BdiscmpCmpnm  string `json:"bdiscmp_cmpnm"`   // 사채권 발행회사(회사명)
	BdiscmpNt     string `json:"bdiscmp_nt"`      // 사채권 발행회사(국적)
	BdiscmpRp     string `json:"bdiscmp_rp"`      // 사채권 발행회사(대표자)
	BdiscmpCpt    string `json:"bdiscmp_cpt"`     // 사채권 발행회사(자본금(원))
	BdiscmpRlCmpn string `json:"bdiscmp_rl_cmpn"` // 사채권 발행회사(회사와 관계)
	BdiscmpTisstk string `json:"bdiscmp_tisstk"`  // 사채권 발행회사(발행주식 총수(주))
	BdiscmpMbsn   string `json:"bdiscmp_mbsn"`    // 사채권 발행회사(주요사업)
	TrfdtlBdFta   string `json:"trfdtl_bd_fta"`   // 양도내역(사채의 권면(전자등록)총액(원))
	TrfdtlTrfprc  string `json:"trfdtl_trfprc"`   // 양도내역(양도금액(원)(A))
	TrfdtlTast    string `json:"trfdtl_tast"`     // 양도내역(총자산(원)(B))
	TrfdtlTastVs  string `json:"trfdtl_tast_vs"`  // 양도내역(총자산대비(%)(A/B))
	TrfdtlEcpt    string `json:"trfdtl_ecpt"`     // 양도내역(자기자본(원)(C))
	TrfdtlEcptVs  string `json:"trfdtl_ecpt_vs"`  // 양도내역(자기자본대비(%)(A/C))
	TrfPp         string `json:"trf_pp"`          // 양도목적
	TrfPrd        string `json:"trf_prd"`         // 양도예정일자
	DlptnCmpnm    string `json:"dlptn_cmpnm"`     // 거래상대방(회사명/성명)
	DlptnCpt      string `json:"dlptn_cpt"`       // 거래상대방(자본금(원))
	DlptnMbsn     string `json:"dlptn_mbsn"`      // 거래상대방(주요사업)
	DlptnHoadd    string `json:"dlptn_hoadd"`     // 거래상대방(본점소재지(주소))
	DlptnRlCmpn   string `json:"dlptn_rl_cmpn"`   // 거래상대방(회사와의 관계)
	DlPym         string `json:"dl_pym"`          // 거래대금지급
	ExevlAtn      string `json:"exevl_atn"`       // 외부평가에 관한 사항(외부평가 여부)
	ExevlBsRs     string `json:"exevl_bs_rs"`     // 외부평가에 관한 사항(근거 및 사유)
	ExevlIntn     string `json:"exevl_intn"`      // 외부평가에 관한 사항(외부평가기관의 명칭)
	ExevlPd       string `json:"exevl_pd"`        // 외부평가에 관한 사항(외부평가 기간)
	ExevlOp       string `json:"exevl_op"`        // 외부평가에 관한 사항(외부평가 의견)
	Bddd          string `json:"bddd"`            // 이사회결의일(결정일)
	OdAAtT        string `json:"od_a_at_t"`       // 사외이사 참석여부(참석(명))
	OdAAtB        string `json:"od_a_at_b"`       // 사외이사 참석여부(불참(명))
	AdtAAtn       string `json:"adt_a_atn"`       // 감사(사외이사가 아닌 감사위원) 참석여부
	FtcSttAtn     string `json:"ftc_stt_atn"`     // 공정거래위원회 신고대상 여부
	PoptCtrAtn    string `json:"popt_ctr_atn"`    // 풋옵션 등 계약 체결여부
	PoptCtrCn     string `json:"popt_ctr_cn"`     // 계약내용
}
```

각 메서드는 위 패턴으로 작성한다.

## 에러 처리

기존 재사용: 데이터 없음 → `opendart.ErrNoData`, 그 외 status → `*opendart.APIError`.

## 테스트 전략

- `material/transfer_test.go`: 기존 `material/client_test.go` 의 `newTestClient` 재사용
  (route map 값은 bare 파일명).
- 8개 메서드 각각 fixture 디코딩 → 대표 필드 검증(머리 + 타입 고유 필드 + 거래상대방/외부평가 1 + 거버넌스).
- fixture 는 실 API 캡처 권장(불가 시 docs 스키마 일치 샘플).
- `integration_test.go` 에 통합 케이스 1~2개(`//go:build integration`, ErrNoData skip 허용).

## 컨벤션 (기존 유지)

- 모든 item struct 필드에 한글 코멘트, 도메인 주석 한국어.
- 표준 net/http(httpclient 재사용), 응답 캐싱 없음, string 유지, UTF-8.
- README "커버리지" DS005 줄에 "양수도(영업/유형자산/타법인주식/주권사채권 양수·양도)" 추가.

## 비범위 (후속 plan)

- 양수도 Sub-2: 자산양수도(기타)·풋백옵션(`astInhtrfEtcPtbkOpt`, 6필드) + 주식교환·이전(`stkExtrDecsn`, 56필드).
- DS005 합병·분할(3) / 해외상장(4). DS006 증권신고서(6). DS002 개인별 보수 Ver 2.0 2종.
