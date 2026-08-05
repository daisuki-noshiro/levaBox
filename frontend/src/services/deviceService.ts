const SIMULATED_BATTERY_LEVEL = 78

/**
 * 设备状态接口占位。
 * 后续在此函数中调用 Wails 生成的 Go 绑定，读取 Windows 实际电量。
 */
export async function getBatteryLevel(): Promise<number> {
  return SIMULATED_BATTERY_LEVEL
}
