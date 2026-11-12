import { useState, useMemo } from 'react'
import { useSalesReport } from '@/hooks/useSalesReport'
import { SectionCards, type StatCardData } from '@/components/manage/report/section-cards'
import { Card, CardContent, CardDescription, CardHeader, CardTitle, CardAction } from '@/components/ui/card'
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select'
import { Loader2, CalendarCheck, Users, Wallet, DollarSign } from 'lucide-react'
import { Area, AreaChart, CartesianGrid, XAxis } from "recharts"
import {
  ChartContainer,
  ChartTooltip,
  ChartTooltipContent,
} from "@/components/ui/chart"
import {
  ToggleGroup,
  ToggleGroupItem,
} from "@/components/ui/toggle-group"
import { bookingStatusChartConfig } from '@/utils/chartConfig'

export default function SaleDashboardPage() {

  const currentDate = useMemo(() => new Date(), [])
  const [selectedYear, setSelectedYear] = useState(currentDate.getFullYear())
  const [selectedMonth, setSelectedMonth] = useState(currentDate.getMonth() + 1)
  const [timeRange, setTimeRange] = useState("30d")

  const { data, isLoading, error, refetch } = useSalesReport({
    year: selectedYear,
    month: selectedMonth
  })

  // Generate year options (last 5 years)
  const yearOptions = useMemo(() => {
    const years = []
    const currentYear = currentDate.getFullYear()
    for (let i = 0; i < 5; i++) {
      years.push(currentYear - i)
    }
    return years
  }, [currentDate])

  // Generate month options
  const monthOptions = [
    { value: 1, label: 'มกราคม' },
    { value: 2, label: 'กุมภาพันธ์' },
    { value: 3, label: 'มีนาคม' },
    { value: 4, label: 'เมษายน' },
    { value: 5, label: 'พฤษภาคม' },
    { value: 6, label: 'มิถุนายน' },
    { value: 7, label: 'กรกฎาคม' },
    { value: 8, label: 'สิงหาคม' },
    { value: 9, label: 'กันยายน' },
    { value: 10, label: 'ตุลาคม' },
    { value: 11, label: 'พฤศจิกายน' },
    { value: 12, label: 'ธันวาคม' },
  ]

  // Filter data based on time range
  const filteredData = useMemo(() => {
    if (!data?.daily_status_report) return []

    const allData = data.daily_status_report
    if (timeRange === "all") return allData

    const daysToShow = timeRange === "7d" ? 7 : 30
    return allData.slice(-daysToShow)
  }, [data, timeRange])

  // Prepare stats data
  const statsData: StatCardData[] = useMemo(() => {
    if (!data) return []

    return [
      {
        title: "จำนวนการจองทั้งหมด",
        value: data.total_bookings.toLocaleString(),
        icon: CalendarCheck,
        iconColor: "text-muted-foreground",
      },
      {
        title: "จำนวนที่นั่งทั้งหมด",
        value: data.total_seats.toLocaleString(),
        icon: Users,
        iconColor: "text-muted-foreground",
      },
      {
        title: "ยอดขายรวม",
        value: `฿${data.total_sales.toLocaleString()}`,
        icon: Wallet,
        iconColor: "text-muted-foreground",
      },
      {
        title: "ยอดชำระที่เก็บได้",
        value: `฿${data.actual_revenue.toLocaleString()}`,
        icon: DollarSign,
        iconColor: "text-muted-foreground",
      },
    ]
  }, [data])

  return (
    <div className="space-y-6">
      <div className="grid grid-cols-1 gap-4 md:flex items-center justify-between ">
        <h1 className="text-3xl font-bold">แดชบอร์ด</h1>
        <div className="flex gap-2">
          <Select value={selectedMonth.toString()} onValueChange={(value) => setSelectedMonth(parseInt(value))}>
            <SelectTrigger className="w-[140px]">
              <SelectValue placeholder="เลือกเดือน" />
            </SelectTrigger>
            <SelectContent>
              {monthOptions.map((month) => (
                <SelectItem key={month.value} value={month.value.toString()}>
                  {month.label}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
          <Select value={selectedYear.toString()} onValueChange={(value) => setSelectedYear(parseInt(value))}>
            <SelectTrigger className="w-[120px]">
              <SelectValue placeholder="เลือกปี" />
            </SelectTrigger>
            <SelectContent>
              {yearOptions.map((year) => (
                <SelectItem key={year} value={year.toString()}>
                  {year}
                </SelectItem>
              ))}
            </SelectContent>
          </Select>
        </div>
      </div>

      {isLoading ? (
        <div className="flex items-center justify-center py-12">
          <Loader2 className="h-8 w-8 animate-spin" />
        </div>
      ) : error ? (
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-red-500">เกิดข้อผิดพลาดในการโหลดข้อมูล</p>
            <button onClick={() => refetch()} className="mt-4 px-4 py-2 bg-primary text-white rounded-md">
              ลองอีกครั้ง
            </button>
          </CardContent>
        </Card>
      ) : data ? (
        <>
          <SectionCards stats={statsData} />

          <Card className="@container/card">
            <CardHeader>
              <CardTitle>รายงานสถานะการจองรายวัน</CardTitle>
              <CardDescription>
                <span className="hidden @[540px]/card:block">
                  ข้อมูลสถานะการจองประจำเดือน {monthOptions.find(m => m.value === selectedMonth)?.label} {selectedYear}
                </span>
                <span className="@[540px]/card:hidden">
                  {monthOptions.find(m => m.value === selectedMonth)?.label} {selectedYear}
                </span>
              </CardDescription>
              <CardAction>
                <ToggleGroup
                  type="single"
                  value={timeRange}
                  onValueChange={setTimeRange}
                  variant="outline"
                  className="hidden *:data-[slot=toggle-group-item]:!px-4 @[767px]/card:flex"
                >
                  <ToggleGroupItem value="all">ทั้งหมด</ToggleGroupItem>
                  <ToggleGroupItem value="30d">30 วันล่าสุด</ToggleGroupItem>
                  <ToggleGroupItem value="7d">7 วันล่าสุด</ToggleGroupItem>
                </ToggleGroup>
                <Select value={timeRange} onValueChange={setTimeRange}>
                  <SelectTrigger
                    className="flex w-40 **:data-[slot=select-value]:block **:data-[slot=select-value]:truncate @[767px]/card:hidden"
                    size="sm"
                    aria-label="Select a value"
                  >
                    <SelectValue placeholder="ทั้งหมด" />
                  </SelectTrigger>
                  <SelectContent className="rounded-xl">
                    <SelectItem value="all" className="rounded-lg">
                      ทั้งหมด
                    </SelectItem>
                    <SelectItem value="30d" className="rounded-lg">
                      30 วันล่าสุด
                    </SelectItem>
                    <SelectItem value="7d" className="rounded-lg">
                      7 วันล่าสุด
                    </SelectItem>
                  </SelectContent>
                </Select>
              </CardAction>
            </CardHeader>
            <CardContent className="px-2 pt-4 sm:px-6 sm:pt-6">
              <ChartContainer
                config={bookingStatusChartConfig}
                className="aspect-auto h-[400px] w-full"
              >
                <AreaChart data={filteredData}>
                  <defs>
                    <linearGradient id="fillPending" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="var(--color-pending)" stopOpacity={0.8} />
                      <stop offset="95%" stopColor="var(--color-pending)" stopOpacity={0.1} />
                    </linearGradient>
                    <linearGradient id="fillProcessing" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="var(--color-processing)" stopOpacity={0.8} />
                      <stop offset="95%" stopColor="var(--color-processing)" stopOpacity={0.1} />
                    </linearGradient>
                    <linearGradient id="fillConfirm" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="var(--color-confirm)" stopOpacity={0.8} />
                      <stop offset="95%" stopColor="var(--color-confirm)" stopOpacity={0.1} />
                    </linearGradient>
                    <linearGradient id="fillCancel" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="var(--color-cancel)" stopOpacity={0.8} />
                      <stop offset="95%" stopColor="var(--color-cancel)" stopOpacity={0.1} />
                    </linearGradient>
                    <linearGradient id="fillCompleted" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="var(--color-completed)" stopOpacity={0.8} />
                      <stop offset="95%" stopColor="var(--color-completed)" stopOpacity={0.1} />
                    </linearGradient>
                    <linearGradient id="fillExpired" x1="0" y1="0" x2="0" y2="1">
                      <stop offset="5%" stopColor="var(--color-expired)" stopOpacity={0.8} />
                      <stop offset="95%" stopColor="var(--color-expired)" stopOpacity={0.1} />
                    </linearGradient>
                  </defs>
                  <CartesianGrid vertical={false} />
                  <XAxis
                    dataKey="date"
                    tickLine={false}
                    axisLine={false}
                    tickMargin={8}
                    minTickGap={32}
                    tickFormatter={(value) => {
                      const date = new Date(value)
                      return date.toLocaleDateString("th-TH", {
                        month: "short",
                        day: "numeric",
                      })
                    }}
                  />
                  <ChartTooltip
                    cursor={false}
                    content={
                      <ChartTooltipContent
                        labelFormatter={(value) => {
                          return new Date(value).toLocaleDateString("th-TH", {
                            year: "numeric",
                            month: "long",
                            day: "numeric",
                          })
                        }}
                        indicator="dot"
                      />
                    }
                  />
                  <Area
                    dataKey="expired"
                    type="monotone"
                    fill="url(#fillExpired)"
                    stroke="var(--color-expired)"
                    stackId="a"
                  />
                  <Area
                    dataKey="completed"
                    type="monotone"
                    fill="url(#fillCompleted)"
                    stroke="var(--color-completed)"
                    stackId="a"
                  />
                  <Area
                    dataKey="cancel"
                    type="monotone"
                    fill="url(#fillCancel)"
                    stroke="var(--color-cancel)"
                    stackId="a"
                  />
                  <Area
                    dataKey="confirm"
                    type="monotone"
                    fill="url(#fillConfirm)"
                    stroke="var(--color-confirm)"
                    stackId="a"
                  />
                  <Area
                    dataKey="processing"
                    type="monotone"
                    fill="url(#fillProcessing)"
                    stroke="var(--color-processing)"
                    stackId="a"
                  />
                  <Area
                    dataKey="pending"
                    type="monotone"
                    fill="url(#fillPending)"
                    stroke="var(--color-pending)"
                    stackId="a"
                  />
                </AreaChart>
              </ChartContainer>
            </CardContent>
          </Card>
        </>
      ) : (
        <Card>
          <CardContent className="py-12 text-center">
            <p className="text-muted-foreground">ไม่มีข้อมูล</p>
          </CardContent>
        </Card>
      )}
    </div>
  )
}