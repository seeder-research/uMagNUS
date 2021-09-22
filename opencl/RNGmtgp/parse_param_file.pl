#!/usr/bin/perl

use strict;
use warnings;
use Getopt::Long qw(GetOptions);
use File::Temp qw/ tempfile tempdir /;
use File::Copy;

# Initiaize globally required variables
my $infile = '';
my $ofile = '';
my $chkHeader = 0;
my $ioffset = 0;
my $total_num = 0;
my $debug = 0;
my $mexpVal = 0;
my $ofh;
my $tmpfname;
my $initString="    {\n";

# Parse program switches
GetOptions(
    'input|i=s' => \$infile,
    'output|o=s' => \$ofile,
    'header=i' => \$chkHeader,
    'input_offset|io=i' => \$ioffset,
    'count|c=i' => \$total_num,
    'debug|d=i' => \$debug,
) or die "Usage: $0 --input|i <filename1> --output|o <filename2> --header <int>\n";

# Error check for unknown input file name
if ($infile eq "") {
	die "Usage: $0 --input|i <filename1> --output|o <filename2> --header <int>\n";
}

# Error check for opening input file
open(my $ifh, '<:encoding(UTF-8)', $infile)
  or die "Could not open file '$infile' $!\n";

# Determine whether output to stdout or output file
# Strategy for output file is to ALWAYS write to a temp file and copy over
# to user defined output file at the end
if ($ofile) {
	($ofh, $tmpfname) = tempfile('tempXXXXXXXXXX', UNLINK => 1);
} else {
	$ofh = *STDOUT;
}

# $row should always point to line of file we are processing...
# Read first row of file first
my $row = <$ifh>;

# If first row contains headers, read next row
$row = <$ifh> if ($chkHeader);

# If we want to discard some rows in the beginning...
if ($ioffset > 0) {
	for (my $idx1 = 0; $idx1 < $ioffset - 1; $idx1++) {
		$row = <$ifh>;
	}
}

# Process the rows we care about and count them. We only process
# up to $total_num of rows.
my $actual_count = 0;
for (my $idx1 = 0; $idx1 < $total_num; $idx1++) {
    # Drop the ending newline
	chomp $row;

	# Separate the comma-separated line into the entries
	my @row_elem_arr = split(',', $row);

	# Error check if the number of entries is unexpected
	my $num_elem = scalar(@row_elem_arr);
	if ($num_elem < 66) {
	    print "Detected invalid line. Exitiing...\n";
		last;
	}

	# Line is valid and so we increment the counter
	$actual_count++;

	# Grab the SHA1 entry...
	my $sha1 = $row_elem_arr[0];
	my $sha_str = "        [21]string{";
	my @sha_str_arr = split('', $sha1);

	# Error check the SHA1 entry
	if (scalar(@sha_str_arr) != 42) {
	    print "Detected invalid sha1 key. Exiting...\n";
		last;
	}

	# SHA1 entry is valid and we process into the string we should output
	my $sha_tmp_str = '';
	for (my $idx2 = 1; $idx2 < 41; $idx2++) {
	    if ($idx2 % 2 != 0) {
			$sha_tmp_str = $sha_str_arr[$idx2];
		} else {
			$sha_tmp_str = $sha_tmp_str . $sha_str_arr[$idx2];
			$sha_str = $sha_str . '"0x' . $sha_tmp_str . '"';
			if ($idx2 == 40) {
				$sha_str = $sha_str . ',"0x00"}';
			} else {
				$sha_str = $sha_str . ',';
			}
		}
	}

	# Grab the mexp entry and process into string we should output
	$mexpVal = $row_elem_arr[1];
	my $mexp = "        int(" . $row_elem_arr[1] . "),\n";

	# Grab the index of the RNG parameter
	my $rid = $row_elem_arr[3];

	# Grab the pos entry and process into string we should output
	my $pos = "        int(" . $row_elem_arr[4] . "),\n";

	# Grab the sh1 entry and process into string we should output
	my $sh1 = "        int(" . $row_elem_arr[5] . "),\n";

	# Grab the sh2 entry and process into string we should output
	my $sh2 = "        int(" . $row_elem_arr[6] . "),\n";

	# Grab the mask entry and process into string we should output
	my $mask = "        uint32(" . $row_elem_arr[15] . "),\n";

	# Grab the weight
	my $weight = $row_elem_arr[16];

	# Grab the delta
	my $delta = $row_elem_arr[17];

	# Process comment line from RNG index, weight and delta
	my $comment = "        /* No." . $rid . " delta:" . $delta . " weight:" . $weight . " */\n";

	# Process RNG tbl and convert to output string
	my $tbl_string="        [16]uint32{";
	for (my $idx2 = 0; $idx2 < 16; $idx2++) {
		$tbl_string = $tbl_string . $row_elem_arr[18+$idx2];
		if ($idx2 == 15) {
			$tbl_string = $tbl_string . "},\n";
		} else {
			$tbl_string = $tbl_string . ",";
		}
	}

	# Process RNG tmp and convert to output string
	my $temper_string="        [16]uint32{";
	for (my $idx2 = 0; $idx2 < 16; $idx2++) {
		$temper_string = $temper_string . $row_elem_arr[34+$idx2];
		if ($idx2 == 15) {
			$temper_string = $temper_string . "},\n";
		} else {
			$temper_string = $temper_string . ",";
		}
	}

	# Process RNG flt_tmp_tbl and convert to output string
	my $flt_temper_tbl="        [16]uint32{";
	for (my $idx2 = 0; $idx2 < 16; $idx2++) {
		$flt_temper_tbl = $flt_temper_tbl . $row_elem_arr[50+$idx2];
		if ($idx2 == 15) {
			$flt_temper_tbl = $flt_temper_tbl . "},\n";
		} else {
			$flt_temper_tbl = $flt_temper_tbl . ",";
		}
	}

	# Print out the detected row entries for debugging
	if ($debug) {
		print "\n";
		for (my $idx = 0; $idx < $num_elem; $idx++) {
			print "$row_elem_arr[$idx]";
			if ($idx < $num_elem - 1) {
				print(",");
			} else {
				print("\n");
			}
		}
	}

	# Build the output string corresponding to detected line parameters
	my $dataString = $initString . $comment . $mexp . $pos . $sh1 . $sh2 . $tbl_string . $temper_string . $flt_temper_tbl . $mask . $sha_str . "     },";

	# If this is the first line, we need to print the
	# leading characters in the full output code. Otherwise,
	# need to print the characters between definitions of
	# RNG parameters
	if ($idx1 == 0) {
	    print $ofh "var MTGP32_params_fast_" . $mexpVal . " = []MTGP32dc_params_fast_t{\n" . $dataString;
	} else {
		print $ofh "\n" . $dataString;
	}

	# Current $row has been processed so we update
	# $row by moving to next line in file
	$row = <$ifh>;
}

# We are done grabbing required RNG parameters. Close input file
close $ifh;

# Print trailing characters of output string
print $ofh "\n};\n";
print $ofh "const mtgpdc_params_" . $mexpVal . "_num = " . $actual_count . "\n";

# If output file was given, we need to close the handle to the temporary
# output file and copy it to the user defined file
if ($ofile) {
	close $ofh;
	copy($tmpfname,$ofile);
}
