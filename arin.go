package main

import "encoding/xml"

type ARINroot struct {
	ARINnet *ARINnet `xml:"net,omitempty" json:"net,omitempty"`
}

type ARINnet struct {
	ARINregistrationDate *ARINregistrationDate `xml:"registrationDate,omitempty"json:"registrationDate,omitempty"`
	ARINcustomerRef      *ARINcustomerRef      `xml:"customerRef,omitempty" json:"customerRef,omitempty"`
	ARINendAddress       *ARINendAddress       `xml:"endAddress,omitempty" json:"endAddress,omitempty"`
	ARINhandle           *ARINhandle           `xml:"handle,omitempty" json:"handle,omitempty"`
	ARINname             *ARINname             `xml:"name,omitempty" json:"name,omitempty"`
	ARINnetBlocks        *ARINnetBlocks        `xml:"netBlocks,omitempty" json:"netBlocks,omitempty"`
	ARINstartAddress     *ARINstartAddress     `xml:"startAddress,omitempty" json:"startAddress,omitempty"`
	ARINupdateDate       *ARINupdateDate       `xml:"updateDate,omitempty" json:"updateDate,omitempty"`
}

type ARINregistrationDate struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"registrationDate,omitempty" json:"registrationDate,omitempty"`
}

type ARINcustomerRef struct {
	Attr_handle string   `xml:" handle,attr"  json:",omitempty"`
	Attr_name   string   `xml:" name,attr"  json:",omitempty"`
	Text        string   `xml:",chardata" json:",omitempty"`
	XMLName     xml.Name `xml:"customerRef,omitempty" json:"customerRef,omitempty"`
}

type ARINendAddress struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"endAddress,omitempty" json:"endAddress,omitempty"`
}

type ARINhandle struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"handle,omitempty" json:"handle,omitempty"`
}

type ARINname struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"name,omitempty" json:"name,omitempty"`
}

type ARINnetBlocks struct {
	ARINnetBlock *ARINnetBlock `xml:"netBlock,omitempty" json:"netBlock,omitempty"`
	XMLName      xml.Name      `xml:"netBlocks,omitempty" json:"netBlocks,omitempty"`
}

type ARINnetBlock struct {
	ARINcidrLength   *ARINcidrLength   `xml:"cidrLength,omitempty" json:"cidrLength,omitempty"`
	ARINdescription  *ARINdescription  `xml:"description,omitempty" json:"description,omitempty"`
	ARINendAddress   *ARINendAddress   `xml:"endAddress,omitempty" json:"endAddress,omitempty"`
	ARINstartAddress *ARINstartAddress `xml:"startAddress,omitempty" json:"startAddress,omitempty"`
	ARINtype         *ARINtype         `xml:"type,omitempty" json:"type,omitempty"`
	XMLName          xml.Name          `xml:"netBlock,omitempty" json:"netBlock,omitempty"`
}

type ARINcidrLength struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"cidrLength,omitempty" json:"cidrLength,omitempty"`
}

type ARINdescription struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"description,omitempty" json:"description,omitempty"`
}

type ARINtype struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"type,omitempty" json:"type,omitempty"`
}

type ARINstartAddress struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"startAddress,omitempty" json:"startAddress,omitempty"`
}

type ARINresources struct {
	Attr_inaccuracyReportUrl string             `xml:" inaccuracyReportUrl,attr"  json:",omitempty"`
	Attr_termsOfUse          string             `xml:" termsOfUse,attr"  json:",omitempty"`
	ARINlimitExceeded        *ARINlimitExceeded `xml:"limitExceeded,omitempty" json:"limitExceeded,omitempty"`
	XMLName                  xml.Name           `xml:"resources,omitempty" json:"resources,omitempty"`
}

type ARINlimitExceeded struct {
	Attr_limit string   `xml:" limit,attr"  json:",omitempty"`
	Text       string   `xml:",chardata" json:",omitempty"`
	XMLName    xml.Name `xml:"limitExceeded,omitempty" json:"limitExceeded,omitempty"`
}

type ARINparentNetRef struct {
	Attr_handle string   `xml:" handle,attr"  json:",omitempty"`
	Attr_name   string   `xml:" name,attr"  json:",omitempty"`
	Text        string   `xml:",chardata" json:",omitempty"`
	XMLName     xml.Name `xml:"parentNetRef,omitempty" json:"parentNetRef,omitempty"`
}

type ARINupdateDate struct {
	Text    string   `xml:",chardata" json:",omitempty"`
	XMLName xml.Name `xml:"updateDate,omitempty" json:"updateDate,omitempty"`
}
