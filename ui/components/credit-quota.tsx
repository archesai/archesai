"use client";
import { useOrganizationsControllerFindOne } from "@/generated/archesApiComponents";
import { useSidebar } from "@/hooks/useSidebar";

export const CreditQuota = () => {
  const { isCollapsed } = useSidebar();
  const { defaultOrgname } = useAuth();
  const { data: organization } = useOrganizationsControllerFindOne({
    pathParams: {
      orgname: defaultOrgname,
    },
  });

  if (isCollapsed) {
    return (
      <CreditCircularChart
        remaining={organization?.credits || 0}
        total={60000}
      />
    );
  }
  return (
    <div className="bg-muted inter p-3 w-full rounded-lg text-xs flex flex-col gap-3.5 ">
      <div className="flex justify-between text-gray-alpha-500">
        <div>Credit quota</div>
        <div>
          <Link className="outline-black" href="/settings/organization/billing">
            <span className="text-gray-alpha-900 font-medium inter text-xs">
              Upgrade
            </span>
          </Link>
        </div>
      </div>
      <div className="font-medium flex items-center gap-2.5">
        <div className="stack gap-0.5 flex-grow">
          <div className="inter flex justify-between">
            <div>Total</div>
            <div className="tabular-nums">{organization?.credits}</div>
          </div>
          <div className="inter flex justify-between">
            <div>Remaining</div>
            <div className="tabular-nums">{organization?.credits}</div>
          </div>
        </div>
        <div>
          <CreditCircularChart
            remaining={organization?.credits || 0}
            total={60000}
          />
        </div>
      </div>
    </div>
  );
};

import { ChartConfig, ChartContainer } from "@/components/ui/chart";
import { useAuth } from "@/hooks/useAuth";
import Link from "next/link";
import {
  PolarGrid,
  PolarRadiusAxis,
  RadialBar,
  RadialBarChart,
} from "recharts";

export const description = "A radial chart with a custom shape";

const chartData = [
  { browser: "safari", fill: "var(--color-safari)", visitors: 1260 },
];

const chartConfig = {
  safari: {
    color: "hsl(var(--chart-2))",
    label: "Safari",
  },
  visitors: {
    label: "Visitors",
  },
} satisfies ChartConfig;

export function CreditCircularChart({
  remaining,
  total,
}: {
  remaining: number;
  total: number;
}) {
  return (
    <ChartContainer
      className="mx-auto aspect-square h-[25px]"
      config={chartConfig}
    >
      <RadialBarChart
        data={chartData}
        endAngle={(remaining / total) * 360}
        height={200}
        innerRadius={10}
        outerRadius={20}
      >
        <PolarGrid
          gridType="circle"
          polarRadius={[86, 74]}
          radialLines={false}
          stroke="none"
        />
        <RadialBar background dataKey="visitors" />
        <PolarRadiusAxis
          axisLine={false}
          tick={false}
          tickLine={false}
        ></PolarRadiusAxis>
      </RadialBarChart>
    </ChartContainer>
  );
}
