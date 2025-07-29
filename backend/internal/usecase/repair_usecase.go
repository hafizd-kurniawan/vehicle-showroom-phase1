// This file would contain the business logic for repairs.
// Key logic points:
// - When a repair is created, the associated vehicle's status is set to 'in_repair'.
// - When adding a part, it checks stock, updates the spare part quantity,
//   updates the repair's total cost, and logs a stock movement.
// - When a repair is marked 'completed', it calculates the final total_repair_cost
//   on the vehicle and updates the vehicle's status to 'ready_to_sell'.
// - It handles all the complex interactions between repairs, vehicles, and spare parts.
