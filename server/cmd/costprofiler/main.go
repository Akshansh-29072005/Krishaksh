package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
)

type Input struct {
	ScansPerMonth           int     `json:"scans_per_month"`
	MAU                     int     `json:"mau"`
	AvgImagesMB             float64 `json:"avg_images_mb"`
	GeminiCostPerScanUSD    float64 `json:"gemini_cost_per_scan_usd"`
	OpenAICostPerScanUSD    float64 `json:"openai_cost_per_scan_usd"`
	FallbackRate            float64 `json:"fallback_rate"`
	S3StorageCostPerGBUSD   float64 `json:"s3_storage_cost_per_gb_usd"`
	QueueInfraMonthlyUSD    float64 `json:"queue_infra_monthly_usd"`
	RDSMonthlyUSD           float64 `json:"rds_monthly_usd"`
	EC2MonthlyUSD           float64 `json:"ec2_monthly_usd"`
	ObservabilityMonthlyUSD float64 `json:"observability_monthly_usd"`
}

func main() {
	inFile := flag.String("in", "", "input json")
	flag.Parse()
	if *inFile == "" {
		fmt.Println("usage: costprofiler -in cost-input.json")
		os.Exit(1)
	}
	b, err := os.ReadFile(*inFile)
	if err != nil {
		panic(err)
	}
	var in Input
	if err := json.Unmarshal(b, &in); err != nil {
		panic(err)
	}
	ai := float64(in.ScansPerMonth) * (in.GeminiCostPerScanUSD*(1-in.FallbackRate) + in.OpenAICostPerScanUSD*in.FallbackRate)
	s3gb := (float64(in.ScansPerMonth) * in.AvgImagesMB) / 1024.0
	s3 := s3gb * in.S3StorageCostPerGBUSD
	infra := in.QueueInfraMonthlyUSD + in.RDSMonthlyUSD + in.EC2MonthlyUSD + in.ObservabilityMonthlyUSD
	total := ai + s3 + infra
	costPerMAU := 0.0
	if in.MAU > 0 {
		costPerMAU = total / float64(in.MAU)
	}
	fmt.Printf("ai_monthly_usd=%.2f\n", ai)
	fmt.Printf("s3_monthly_usd=%.2f\n", s3)
	fmt.Printf("infra_monthly_usd=%.2f\n", infra)
	fmt.Printf("total_monthly_usd=%.2f\n", total)
	fmt.Printf("cost_per_mau_usd=%.4f\n", costPerMAU)
}
