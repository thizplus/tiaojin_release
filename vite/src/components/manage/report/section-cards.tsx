
import {
  Card,
  CardDescription,
  CardHeader,
  CardTitle,
} from "@/components/ui/card"
import type { LucideIcon } from "lucide-react"

export interface StatCardData {
  title: string
  value: string | number
  icon: LucideIcon
  iconColor: string
}

interface SectionCardsProps {
  stats: StatCardData[]
}

export function SectionCards({ stats }: SectionCardsProps) {
  return (
    <div className="grid grid-cols-1 gap-4 md:grid-cols-2 lg:grid-cols-4">
      {stats.map((stat, index) => {
        const Icon = stat.icon
        return (
          <Card key={index} className="@container/card">
            <CardHeader>


<div className="flex items-center gap-2">
<Icon className={`size-12 ${stat.iconColor}`} />
<div className=" items-center ">
                <CardDescription className="text-sm font-medium">{stat.title}</CardDescription>
                <CardTitle className="text-2xl font-semibold tabular-nums @[250px]/card:text-3xl">
                {stat.value}
              </CardTitle>
              </div>
</div>



            
            </CardHeader>
          </Card>
        )
      })}
    </div>
  )
}
