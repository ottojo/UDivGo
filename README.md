# UDivGo
Voltage Divider Calculator that presents the optimal Voltage Divider consisting of a given set of availiable resistors.
## Usage
```
UDivGo
  -Uin float
        Input voltage in volts (e.g. 5.0)
  -Uout float
        Output voltag in volts (e.g. 3.3)
  -n int
        Number of dividers to calculate (e.g. 3)
  -resistors string
        File containing availiable resistors: Each line should contain the package description and,
        seperated by at least one space, the resistance in ohms.
        e.g.  SMD_1206 1000
        would describe a SMD 1kΩ resistor. The package is not used.
```
