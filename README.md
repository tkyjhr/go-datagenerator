# go-datagenerator

Generate some sequence of data and output to stdout or file.
`reader.go` defines classes which implements io.Reader to generate data.

```
NAME:
   datagenerator - Output data

USAGE:
   datagenerator (-d [Data]) (-o [OutputPath]) [Size]

   [Data]       : Data to generate. See below.
   [OutputPath] : Output file path. Directries are created if necessary.
                  If ommitted, data is output to stdout.
   [Size]       : Size in bytes. KB and MB can be used as suffix for (* 1024) and (* 1024 * 1024) respectively.
                  Hexadecimal value can be used with 0x prefix.
                  Maximum size is 1GB.


OPTIONS:
   --data, -d   Specify content of output data. 0 is the default.
                0    : Fill with 0x00.
                f    : Fill with 0xFF.
                r    : Fiil with random bytes.
                ra   : Fill with random alphabet character.
                ran  : Fill with random alphabet and numeric characters.
                ctr  : Fill with bytes increasing from 0x00 to 0xFF. After 0xFF it backs to 0x00.
                ctr2 : Fill with uint16(BigEndian) increasing from 0x0000 to 0xFFFF.
                       After 0xFFFF, the value starts from 0x0000 again.
                       The size must be a multiple of 2.
                ctr4 : Fill with uint32(BigEndian) increasing from 0x00000000 to 0xFFFFFFFF.
                       After 0xFFFFFFFF, the value starts from 0x0000 again.
                       The size must be a multiple of 4.
                ctr8 : Fill with uint64(BigEndian) increasing from 0x0000000000000000 to 0xFFFFFFFFFFFFFFFF.
                       After 0xFFFFFFFFFFFFFFFF, the value starts from 0x0000000000000000 again.
                       The size must be a multiple of 8.

   -o           Output file path.
   --help, -h   show help
```
