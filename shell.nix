with import <nixpkgs> { };

let
  nodejs = nodejs-12_x;

  userName = "admin";
  appName = "test-app2";
  gtwName = "test-gtw";
  mkDevName = class:
    let base = "test-dev";
    in if class == "a" then
      "${base}-a"
    else if class == "b" then
      "${base}-b"
    else if class == "c" then
      "${base}-c"
    else
      throw "unknown class '${class}'";

  frequencyPlan = "EU_863_870";
  #frequencyPlan = "US_902_928_FSB_2";

  stackPort = "1885";
  udpPort = "1700";

  pwd = toString ./.;
  cmdDir = "${pwd}/cmd";

  goRun = "${go}/bin/go run";
  cliCmd = let conf = "${pwd}/ttn-lw-cli.yml";
  in ''
    $([[ -d '${cmdDir}/tti-lw-cli' ]] && printf "${goRun} -tags 'tti' ${cmdDir}/tti-lw-cli -c ${conf}" || printf "${goRun} ${cmdDir}/ttn-lw-cli -c ${conf}")'';
  stackCmd = let conf = "${pwd}/ttn-lw-stack.yml";
  in ''
    $([[ -d '${cmdDir}/tti-lw-stack' ]] && printf "${goRun} -tags 'tti' ${cmdDir}/tti-lw-stack -c ${conf}" || printf "${goRun} ${cmdDir}/ttn-lw-stack -c ${conf}")'';
  mageCmd = let mage = "${pwd}/tools/bin/mage";
  in "MAGE=${mage} ${gnumake}/bin/make -C ${pwd} -s ${mage} && pushd ${pwd} && ${mage} init && ${mage}";

  createAPIKeyCmd =
    "${cliCmd} application api-keys create --application-id ${appName} --right-application-all | ${jq}/bin/jq -c -r '.key'";

  makeWaitCmd = event: ''
    read -n 1 -p "Ensure that ${event} and press any key to continue"
  '';

  joinWaitCmd = makeWaitCmd "the device has joined(OTAA) or booted(ABP)";
  uplinkWaitCmd = makeWaitCmd "the device has sent an uplink";

  linkApp = writeShellScriptBin "link-app" ''
    set -xe
    apiKey="$(${createAPIKeyCmd})"
    ${cliCmd} application link set ${appName} --api-key ''${apiKey} ''${@}
    printf "export TTN_API_KEY=\"''${apiKey}\"" > ${pwd}/.envrc.api_key
    ${direnv}/bin/direnv reload
  '';

  bootstrap = writeShellScriptBin "bootstrap" ''
    set -xe
    ${cliCmd} login
    ${cliCmd} application create ${appName} --user-id ${userName}
    ${cliCmd} gateway create ${gtwName} --defaults --user-id ${userName} --frequency_plan_id ${frequencyPlan}
    printf "export BOOTSTRAP_DONE=1" > ${pwd}/.envrc.bootstrap
    ${linkApp}/bin/link-app
  '';

  scheduleClockExample = writeShellScriptBin "schedule-clock-example" ''
    set -xe
    cd ../../github.com/TheThingsNetwork/lorawan-stack-example-clock/controller
    export TTN_APPLICATION_ID=''${TTN_APPLICATION_ID:-${appName}}
    export TTN_DEVICE_ID=''${TTN_DEVICE_ID:-${mkDevName "c"}}
    export TTN_API_KEY=''${TTN_API_KEY:-$(cat ${pwd}/.envrc.api_key)}
    export TTN_SERVER=''${TTN_SERVER:-http://localhost:${stackPort}}
    ${nodejs}/bin/npm run start
  '';

  # TODO: Support class B/C args
  makePushDownlinksCmd = devName:
    lib.concatMapStringsSep "\n" (down: ''
      ${cliCmd} device downlink push ${appName} ${devName} \
                                      --f-port=${toString down.port} \
                                      --frm-payload=${down.payload}
        sleep 1
    '');

  classADownlinks = [
    {
      port = 1;
      payload = "1111";
    }
    {
      port = 42;
      payload = "22";
    }
    {
      port = 120;
      payload = "333333";
    }
  ];

  pushClassADownlinksCmd = devName:
    makePushDownlinksCmd devName classADownlinks;
  pushClassBDownlinksCmd = devName:
    makePushDownlinksCmd devName (classADownlinks ++ [
      # TODO: Add class B downlinks
    ]);
  pushClassCDownlinksCmd = devName:
    makePushDownlinksCmd devName (classADownlinks ++ [
      # TODO: Add class C downlinks
    ]);

  updateIfVersionAtLeast = v1: v2: base: attrs:
    if lib.strings.versionAtLeast v1 v2 then base // attrs else base;

  devices = let
    mkGeneric = mac: phy: {
      "lorawan-phy-version" = phy;
      "lorawan-version" = mac;
    };

    mkGenericOTAA = mac: phy:
      updateIfVersionAtLeast mac "1.1.0" (mkGeneric mac phy // {
        "dev-eui" = "DEADBEEF01020304";
        "join-eui" = "01020304DEADBEEF";
        "resets-join-nonces" = true;
        "root-keys.app-key.key" = "01020304DEADBEEF01020304DEADBEEF";
        "supports-join" = true;
      }) { "root-keys.nwk-key.key" = "DEADBEEF01020304DEADBEEF01020304"; };

    mkGenericABP = mac: phy:
      updateIfVersionAtLeast mac "1.1.0" (mkGeneric mac phy // {
        "supports-join" = false;
        "session.dev-addr" = "DEADBEEF";
        "session.keys.f-nwk-s-int-key.key" = "01020304DEADBEEF01020304DEADBEEF";
        "session.keys.app-s-key.key" = "DEADBEEF01020304DEADBEEF01020304";
      }) {
        "session.keys.s-nwk-s-int-key.key" = "01010203DEADBEEF01010203DEADBEEF";
        "session.keys.nwk-s-enc-key.key" = "01010102DEADBEEF01010102DEADBEEF";
      };

    mkMbedOTAA = mkGenericOTAA;

    mkLMNOTAA = mac: phy:
      updateIfVersionAtLeast mac "1.1.0" (mkGenericOTAA mac phy // {
        "dev-eui" = "DEADBEEF02030405";
        "root-keys.app-key.key" = "02030405DEADBEEF02030405DEADBEEF";
      }) { "root-keys.nwk-key.key" = "DEADBEEF02030405DEADBEEF02030405"; };
  in lib.foldr (ver: acc:
    acc // {
      "generic-${ver.mac}-otaa" = mkGenericOTAA ver.mac ver.phy;
      "generic-${ver.mac}-abp" = mkGenericABP ver.mac ver.phy;
    }) { } [
      {
        mac = "1.0.0";
        phy = "1.0.0";
      }
      {
        mac = "1.0.2";
        phy = "1.0.2-a";
      }
      {
        mac = "1.0.2";
        phy = "1.0.2-b";
      }
      {
        mac = "1.0.3";
        phy = "1.0.3-a";
      }
      {
        mac = "1.1.0";
        phy = "1.1.0-a";
      }
      {
        mac = "1.1.0";
        phy = "1.1.0-b";
      }
    ] // {
      "mbed-1.0.3-otaa" = mkMbedOTAA "1.0.3" "1.0.3-a";
      "lmn-1.0.3-otaa" = mkLMNOTAA "1.0.3" "1.0.3-a";

      "semtech-otaa" = {
        "dev-eui" = "ffffffaa11111110";
        "join-eui" = "0102030405060708";
        "root-keys.app-key.key" = "\${SEMTECH_APPKEY}";
        "lorawan-version" = "1.0.2";
        "lorawan-phy-version" = "1.0.2-a";
        "supports-join" = true;
      };

      "pba-abp" = {
        "session.dev-addr" = "00000001";
        "session.keys.app-s-key.key" = "01020304DEADBEEF01020304DEADBEEF";
        "session.keys.f-nwk-s-int-key.key" = "01020304DEADBEEF01020304DEADBEEF";
        "lorawan-phy-version" = "1.0.3-a";
        "lorawan-version" = "1.0.3";
        "supports-join" = false;
      };

      "uno-otaa" = {
        "dev-eui" = "01020304DEADBEEF";
        "join-eui" = "01020304DEADBEEF";
        "root-keys.app-key.key" = "01020304DEADBEEF01020304DEADBEEF";
        "lorawan-version" = "1.0.2";
        "lorawan-phy-version" = "1.0.2-b";
        "supports-join" = true;
        "resets-join-nonces" = true;
      };
    };

  toFlags = lib.generators.toKeyValue {
    mkKeyValue = k: v: "--${k}=${lib.generators.mkValueStringDefault { } v} \\";
  };

  writeTestDeviceBin = name: devName: flags: script:
    writeShellScriptBin name ''
      [ "''${BOOTSTRAP_DONE}" == "1" ] || ${bootstrap}/bin/bootstrap

      set -xe

      ${cliCmd} device delete ${appName} ${devName} || true

      case "''${1}" in ${
        lib.generators.toKeyValue {
          mkKeyValue = k: v: ''
            "${k}")
            ${cliCmd} device create ${appName} ${devName} \
              --frequency_plan_id=${frequencyPlan} \
              ${toFlags v} ${toFlags flags} ''${@}
            ;;
          '';
        } devices
      }
          *)
          echo "Unknown device ''${1} - specify one of ${
            with builtins;
            toString (attrNames devices)
          }"
          exit 1
          esac
          shift

          ${script}
        '';
in mkShell {
  CGO_ENABLED = "0";
  GO111MODULE = "on";
  GOROOT = "${go}/share/go";
  TTN_LW_GS_UDP_LISTENERS = ":${udpPort}=${frequencyPlan}";

  name = "ttn-env";
  buildInputs = [
    (writeShellScriptBin "mage" ''
      ${mageCmd} ''${@}
    '')

    (writeShellScriptBin "cli" ''
      ${cliCmd} ''${@}
    '')
    (writeShellScriptBin "ttn-lw-cli" ''
      ${cliCmd} ''${@}
    '')
    (writeShellScriptBin "stack" ''
      ${stackCmd} ''${@}
    '')
    (writeShellScriptBin "ttn-lw-stack" ''
      ${stackCmd} ''${@}
    '')

    (writeShellScriptBin "snapcraft" ''
      ${docker}/bin/docker run --rm -i \
        --user $(id -u):$(id -g) \
        --mount type=bind,src=$XDG_CONFIG_HOME/snapcraft,dst=$XDG_CONFIG_HOME/snapcraft \
        --mount type=bind,src=$(pwd),dst=$(pwd) \
        --entrypoint "/snap/bin/snapcraft" \
        -e XDG_CONFIG_HOME \
        -w $(pwd) \
        snapcore/snapcraft:candidate ''${@}
    '')

    (writeShellScriptBin "start-clean-stack" ''
      set -xe
      ${mageCmd} dev:dbStop
      sudo ${coreutils}/bin/rm -rf .env/data
      ${docker-compose}/bin/docker-compose pull redis cockroach
      ${mageCmd} dev:dbStart
      ${mageCmd} dev:initStack
      printf "export BOOTSTRAP_DONE=0" > ${pwd}/.envrc.bootstrap
      ${direnv}/bin/direnv reload
      ${stackCmd} start ''${@}
    '')

    (writeShellScriptBin "create-app" ''
      set -xe
      ${cliCmd} application create ${appName} --user-id ${userName} ''${@}
    '')
    (writeShellScriptBin "delete-app" ''
      set -xe
      ${cliCmd} application delete ${appName} ''${@}
    '')

    bootstrap
    linkApp
    scheduleClockExample

    (let devName = mkDevName "a";
    in writeTestDeviceBin "test-class-a-device" devName { } ''
      ${joinWaitCmd}

      ${makePushDownlinksCmd devName [{
        port = 2;
        payload = "abcd";
      }]}

      ${uplinkWaitCmd}

      ${pushClassADownlinksCmd devName}
    '')

    (let devName = mkDevName "b";
    in writeTestDeviceBin "test-class-b-device" devName {
      supports-class-b = true;
      #"mac-settings.desired-ping-slot-data-rate-index" = "DATA_RATE_0";
    } ''
      ${joinWaitCmd}

      ${pushClassBDownlinksCmd devName}
    '')

    (let devName = mkDevName "c";
    in writeTestDeviceBin "test-class-c-device" devName {
      supports-class-c = true;
    } ''
      ${joinWaitCmd}

      source ${pwd}/.envrc.api_key
      ${scheduleClockExample}/bin/schedule-clock-example

      ${pushClassCDownlinksCmd devName}
    '')
  ] ++ [
    act
    coreutils
    entr
    gnumake
    go
    go-tools
    gopls
    gotools
    jq
    minicom
    mosquitto
    neovim
    nodejs
    perl
    postgresql_12
    protobuf
    rdbtools
    redis
    richgo
    rpm
    travis
    xxd
  ];
}
