import csv
import random

# Tên file output
output_file = "./class.csv"

# Ghi dữ liệu vào CSV
with open(output_file, mode="w", newline="", encoding="utf-8") as file:
    writer = csv.writer(file)

    # Header
    writer.writerow(["CLASS", "SLOT"])

    # Sinh dữ liệu từ 161001 -> 161999
    for i in range(1, 1000):
        class_code = f"161{i:03d}"
        slot = random.randint(40, 70)

        writer.writerow([class_code, slot])

print(f"Đã tạo file {output_file}")