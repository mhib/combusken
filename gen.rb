class Object
  def id
    self
  end
end

for i in %w(rook bishop queen)
  puts <<-HERE
func #{i.id.capitalize}sAttacks(set uint64, occupancy uint64) uint64 {
    res := uint64(0)
    for set > 0 {
      position := bitScan(set)
      set &= ^(1u64 << position)

      res |= #{i.id}Attacks(position, occupancy)
      }
    return res
    }
    HERE
end

