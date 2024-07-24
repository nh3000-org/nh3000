package config

// nh3000 fields
var SIPLoggedOn = false

var SIPCaroot = ""
var SIPClientCert = ""
var SIPClientKey = ""

// sip fields
var SipGeneralCallWaiting = true
var SipGeneralHangup3Way = true
var SipGeneralInCallOSD = true
var SipGeneralCallHistory = 32 // 8 16 32 64 128

var SipSoundCardRingtone = ""
var SipSoundCardSpeaker = ""
var SipSoundCardMicrophone = ""
var SipSoundCardValidate = true
var SipSoundCardOssFragsize = 1024           // 16 32 64 128 256 512 1024
var SipSoundCardAlsaPlayPeriodSize = 1024    // 16 32 64 128 256 512 1024
var SipSoundCardAlsaCapturePeriodSize = 1024 // 16 32 64 128 256 512 1024

var SipRingtonePlay = true
var SipRingtoneTone = ""
var SipRingtoneBackPlay = true
var SipRingToneBackTone = ""

var SipAddressbookLookup = true
var SipAddressbookOverideAddress = true
var SipAddressbookLookupPhoto = true

var SipNetworkPort = 5060
var SipNetworkRtpPort = 8000
var SipNetworkMaxUdpSize = 65535
var SipNetworkMaxTcpSize = 1000000

var SipServerDNS = "192.168.0.15"
var SipServerExpirySeconds = 3600
var SipServerRegisterAtStartup = true
var SipServerAddQvalueToRegistration = "1.000"
var SipServerUseOutboundProxy = false
var SipServerOutboundProxy = ""
var SipServerSendInDialogRequestsProxy = false

var SipVoiceMailboxAddress = 9999
var SipVoiceMailMWIType = "Solicited"
var SipVoiceMailboxUserName = 9999
var SipVoiceMailServer = "192.168.0.15"
var SipVoiceMailServerUseProxy = false
var SipVoiceMailSubscruptionDuration = "3600"

var SipPresenseAvailableAtStartup = true

var SipProtocolCallHoldRFC = "3264" // RFC 2543 3264
var SipProtocolMaxForwardsHeaderMandatory = false
var SipProtocolAllowMissingContactIn200 = true
var SipProtocolPutRegExpiryInHeader = true
var SipProtocolUseCompactHeaderNames = false
var SipProtocolEncodeRouteAsList = true
var SipProtocolUseDomainInHeader = false
var SipProtocolAllowSDPChangeDuringSetup = false
var SipProtocolRedirectionAllow = true
var SipProtocolRedirectAsk = true
var SipProtocolRedirectMax = 1
var SipProtocolExtensionsPRACK = "SUPPORTED" // DISABLED SUPPORTED REQUIRED PREFERRED
var SipProtocolExtensionsReplace = true
var SipProtocolReferAcceptTransfer = true
var SipProtocolReferAsk = true
var SipProtocolReferHoldReferrerWhile = false
var SipProtocolReferHoldReferrerBefore = false
var SipProtocolReferAutoRefresh = false
var SipProtocolReferAttendedAOR = false
var SipProtocolReferAllowTransfer = false
var SipProtocolPrivacyPPreferred = false
var SipProtocolPrivacyPAsserted = false
var SipTransportProtocol = "AUTO"          // TCP UDP
var SipTransportUDOThreshold = "1300bytes" // 1300bytes 1600bytes
var SipTransportNATTransversal = false
var SipTransportNATTransversalUseStaticIP = ""
var SipTransportNATTransversalUseStun = ""
var SipTransportNATKeepAlive = false
var SipTimerNoAnswer = 30
var SipTimerNAtKeepAlive = 30
var SipSecurityEnableRTPEncryption = false

var SipUserAlias []string
var SipUserName []string
var SipUserDomain []string
var SipUserOrganization []string
var SipUserAuthRealm []string
var SipUserAuthName []string
var SipUserAuthPassword []string
var SipUserAkAop []string
var SipUserAkaAmp []string
