"use client"

import * as React from "react"
import { format } from "date-fns"
import { th } from "date-fns/locale"
import { CalendarIcon } from "lucide-react"

import { Button } from "@/components/ui/button"
import { Calendar } from "@/components/ui/calendar"
import { Label } from "@/components/ui/label"
import {
  Popover,
  PopoverContent,
  PopoverTrigger,
} from "@/components/ui/popover"
import { cn } from "@/lib/utils"

interface DatePickerFieldProps {
  id?: string
  label?: string
  placeholder?: string
  value?: Date | string
  onChange?: (date: Date | undefined) => void
  required?: boolean
  error?: string
  disabled?: boolean
  className?: string
}

export function DatePickerField({
  id,
  label,
  placeholder = "เลือกวันที่",
  value,
  onChange,
  required = false,
  error,
  disabled = false,
  className,
}: DatePickerFieldProps) {
  const [open, setOpen] = React.useState(false)

  // Convert string to Date if needed (with Thai timezone handling)
  const dateValue = React.useMemo(() => {
    if (!value) return undefined
    if (value instanceof Date) return value
    if (typeof value === 'string') {
      // Create date in local timezone to avoid UTC conversion issues
      const [year, month, day] = value.split('-').map(Number)
      if (year && month && day) {
        return new Date(year, month - 1, day) // month is 0-indexed
      }
      return undefined
    }
    return undefined
  }, [value])

  const handleSelect = (date: Date | undefined) => {
    if (date) {
      // Format date to YYYY-MM-DD in local timezone
      const year = date.getFullYear()
     // const month = String(date.getMonth() + 1).padStart(2, '0')
    //  const day = String(date.getDate()).padStart(2, '0')
      const formattedDate = new Date(year, date.getMonth(), date.getDate())
      onChange?.(formattedDate)
    } else {
      onChange?.(undefined)
    }
    setOpen(false)
  }

  return (
    <div className={cn("flex flex-col gap-2", className)}>
      {label && (
        <Label htmlFor={id} className="px-1">
          {label}
          {required && <span className="text-destructive ml-1">*</span>}
        </Label>
      )}
      <Popover open={open} onOpenChange={setOpen}>
        <PopoverTrigger asChild>
          <Button
            variant="outline"
            id={id}
            disabled={disabled}
            className={cn(
              "w-full justify-between font-normal",
              !dateValue && "text-muted-foreground",
              error && "border-destructive"
            )}
          >
            {dateValue ? (
              format(dateValue, "dd/MM/yyyy", { locale: th })
            ) : (
              placeholder
            )}
            <CalendarIcon className="h-4 w-4" />
          </Button>
        </PopoverTrigger>
        <PopoverContent className="w-auto overflow-hidden p-0" align="start">
          <Calendar
            mode="single"
            selected={dateValue}
            onSelect={handleSelect}
            captionLayout="dropdown"
            locale={th}
            disabled={disabled}
            startMonth={new Date(2020, 0)}
            endMonth={new Date(2035, 11)}
          />
        </PopoverContent>
      </Popover>
      {error && (
        <p className="text-sm text-destructive px-1">
          {error}
        </p>
      )}
    </div>
  )
}