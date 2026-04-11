package gov

// Curve 投票曲线，用于计算批准/支持阈值
type Curve struct {
	// 曲线类型：LinearDecreasing, SteppedDecreasing, Reciprocal
	Type string

	// LinearDecreasing 参数
	// 从 (0, begin) 开始，线性递减到 (length, end)，然后保持在 end
	LinearBegin  uint32
	LinearEnd    uint32
	LinearLength int64

	// SteppedDecreasing 参数
	// 从 (0, begin) 开始，保持 period 个区块不变，然后下降 step
	StepBegin  uint32
	StepEnd    uint32
	StepSize   uint32
	StepPeriod int64

	// Reciprocal 参数
	// 公式：K/(x/S + x_offset) - y_offset
	ReciprocalFactor  uint32
	ReciprocalXScale  uint32
	ReciprocalXOffset int64
	ReciprocalYOffset int64
}

// CurveType 常量
const (
	CurveTypeLinearDecreasing  = "LinearDecreasing"
	CurveTypeSteppedDecreasing = "SteppedDecreasing"
	CurveTypeReciprocal        = "Reciprocal"
)

// NewLinearDecreasingCurve 创建线性递减曲线
func NewLinearDecreasingCurve(begin, end uint32, length int64) Curve {
	return Curve{
		Type:         CurveTypeLinearDecreasing,
		LinearBegin:  begin,
		LinearEnd:    end,
		LinearLength: length,
	}
}

// NewSteppedDecreasingCurve 创建阶梯递减曲线
func NewSteppedDecreasingCurve(begin, end, step uint32, period int64) Curve {
	return Curve{
		Type:       CurveTypeSteppedDecreasing,
		StepBegin:  begin,
		StepEnd:    end,
		StepSize:   step,
		StepPeriod: period,
	}
}

// NewReciprocalCurve 创建倒数曲线
func NewReciprocalCurve(factor, xScale uint32, xOffset, yOffset int64) Curve {
	return Curve{
		Type:              CurveTypeReciprocal,
		ReciprocalFactor:  factor,
		ReciprocalXScale:  xScale,
		ReciprocalXOffset: xOffset,
		ReciprocalYOffset: yOffset,
	}
}

// Y 计算曲线在给定 X（区块号）处的 Y 值（阈值百分比）
func (c *Curve) Y(x int64) uint32 {
	switch c.Type {
	case CurveTypeLinearDecreasing:
		return c.linearDecreasingY(x)
	case CurveTypeSteppedDecreasing:
		return c.steppedDecreasingY(x)
	case CurveTypeReciprocal:
		return c.reciprocalY(x)
	default:
		return 0
	}
}

func (c *Curve) linearDecreasingY(x int64) uint32 {
	if x >= c.LinearLength {
		return c.LinearEnd
	}
	if c.LinearLength == 0 {
		return c.LinearEnd
	}
	// slope = (begin - end) * x / length
	slope := int64(c.LinearBegin-c.LinearEnd) * x / c.LinearLength
	result := int64(c.LinearBegin) - slope
	if result < int64(c.LinearEnd) {
		return c.LinearEnd
	}
	return uint32(result)
}

func (c *Curve) steppedDecreasingY(x int64) uint32 {
	if c.StepPeriod == 0 || x < c.StepPeriod {
		return c.StepBegin
	}
	numSteps := x / c.StepPeriod
	subValue := numSteps * int64(c.StepSize)

	if subValue >= int64(c.StepBegin) || c.StepBegin-uint32(subValue) <= c.StepEnd {
		return c.StepEnd
	}
	return c.StepBegin - uint32(subValue)
}

func (c *Curve) reciprocalY(x int64) uint32 {
	if c.ReciprocalXScale == 0 {
		return 0
	}
	// K/(x/S + x_offset) - y_offset
	xScale := int64(c.ReciprocalXScale)
	denominator := x/xScale + c.ReciprocalXOffset
	if denominator <= 0 {
		denominator = 1
	}
	result := int64(c.ReciprocalFactor)/denominator - c.ReciprocalYOffset
	if result < 0 {
		return 0
	}
	return uint32(result)
}
