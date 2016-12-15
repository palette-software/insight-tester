# Disable the stupid stuff rpm distros include in the build process by default:
#   Disable any prep shell actions. replace them with simply 'true'
%define __spec_prep_post true
%define __spec_prep_pre true
#   Disable any build shell actions. replace them with simply 'true'
%define __spec_build_post true
%define __spec_build_pre true
#   Disable any install shell actions. replace them with simply 'true'
%define __spec_install_post true
%define __spec_install_pre true
#   Disable any clean shell actions. replace them with simply 'true'
%define __spec_clean_post true
%define __spec_clean_pre true
# Disable checking for unpackaged files ?
#%undefine __check_files

# Use md5 file digest method.
# The first macro is the one used in RPM v4.9.1.1
%define _binary_filedigest_algorithm 1
# This is the macro I find on OSX when Homebrew provides rpmbuild (rpm v5.4.14)
%define _build_binary_file_digest_algo 1

# Use bzip2 payload compression
%define _binary_payload w9.bzdio

# Enable bash specific commands (eg. pushd)
%define _buildshell /bin/bash

Name: palette-insight-sanity-check
Version: %version
Release: %buildrelease
Summary: Palette Insight Sanity Check
AutoReqProv: no
# Seems specifying BuildRoot is required on older rpmbuild (like on CentOS 5)
# fpm passes '--define buildroot ...' on the commandline, so just reuse that.
#BuildRoot: %buildroot
# Add prefix, must not end with / except for root (/)

Prefix: /

Group: default
License: proprietary
Vendor: palette-software.net
URL: http://www.palette-software.com
Packager: Palette Developers <developers@palette-software.com>

# Add the user for the service & setup SELinux
# ============================================

%description
Palette Insight Sanity Check

%define target_install_dir /opt/insight-sanity-check
%prep
# noop

%build
# noop

%install

mkdir -p %{buildroot}/%{target_install_dir}
pushd %{buildroot}/%{target_install_dir}
cp ${GOPATH}/bin/dbcheck .
cp -R %{source_dir}/dbcheck/tests .
cp %{source_dir}/dbcheck/sanity-check.sh .
cp %{source_dir}/dbcheck/Config_template.yml Config.yml
sed -i -e "s/{{ gp_palette_password }}/${GP_PALETTE_PASSWORD}/" ./Config.yml
popd

%files
%defattr(-,insight,insight,-)

# Reject config files already listed or parent directories, then prefix files
# with "/", then make sure paths with spaces are quoted.
%{target_install_dir}/dbcheck
%{target_install_dir}/tests
%{target_install_dir}/sanity-check.sh

# config files can be defined according to this
# http://www-uxsup.csx.cam.ac.uk/~jw35/docs/rpm_config.html
%config %{target_install_dir}/Config.yml

%clean
# noop

%pre
# noop

%post
# noop


%changelog
