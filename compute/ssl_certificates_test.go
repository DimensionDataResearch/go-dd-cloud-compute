package compute

import "testing"

// List SSL domain certificates (successful).
func TestClient_ListSSLDomainCertificatesInNetworkDomain_Success(test *testing.T) {
	testClientRequest(test, &ClientTestConfig{
		Request: func(test *testing.T, client *Client) {
			const networkDomainID = "14dbfacf-0e ec-4964-a0c2-ff3f7390246b"

			page := &Paging{
				PageNumber: 1,
				PageSize:   250,
			}
			domainCertificates, err := client.ListSSLDomainCertificatesInNetworkDomain(networkDomainID, page)
			if err != nil {
				test.Fatal(err)
			}

			verifyListSSLDomainCertificatesInNetworkDomainTestResponse(test, domainCertificates)
		},
		Respond: testRespondOK(listSSLDomainCertificatesInNetworkDomainTestResponse),
	})
}

// Import SSL domain certificate (successful).
func TestClient_ImportSSLDomainCertificate_Success(test *testing.T) {
	expect := expect(test)

	testClientRequest(test, &ClientTestConfig{
		Request: func(test *testing.T, client *Client) {
			sslDomainCertificateID, err := client.ImportSSLDomainCertificate(
				"553f26b6-2a73-42c3-a78b-6116f11291d0",
				"Test-SSL-Domain-Certificate",
				"Test SSL Domain Certificate Description",
				"-----BEGIN CERTIFICATE-----\nMIIFwzC...truncated for documentation brevity purposes...0xRcT7WQYUtUxu\n-----END CERTIFICATE-----",
				"-----BEGIN PRIVATE KEY-----\nBAQEFAASCBKgwggGkl...truncated for documentation brevity purposes...wBPNAUcZVF5umh6GP\n-----END PRIVATE KEY-----",
			)
			if err != nil {
				test.Fatal(err)
			}

			expect.EqualsString("SSLDomainCertificateID", "9e6b496d-5261-4542-91aa-b50c7f569c54", sslDomainCertificateID)
		},
		Respond: testValidateJSONRequestAndRespondOK(importSSLDomainCertificateTestRespose, &importSSLDomainCertificate{}, func(test *testing.T, requestBody interface{}) {
			verifyImportSSLDomainCertificateTestRequest(test, requestBody.(*importSSLDomainCertificate))
		}),
	})
}

/*
 * Test requests.
 */

var importSSLDomainCertificateTestRequest = `
{
    "networkDomainId": "553f26b6-2a73-42c3-a78b-6116f11291d0",
    "name": "Test-SSL-Domain-Certificate",
    "description": "Test SSL Domain Certificate Description",
    "key": "-----BEGIN PRIVATE KEY-----\nBAQEFAASCBKgwggGkl...truncated for documentation brevity purposes...wBPNAUcZVF5umh6GP\n-----END PRIVATE KEY-----",
    "certificate": "-----BEGIN CERTIFICATE-----\nMIIFwzC...truncated for documentation brevity purposes...0xRcT7WQYUtUxu\n-----END CERTIFICATE-----"
}
`

func verifyImportSSLDomainCertificateTestRequest(test *testing.T, request *importSSLDomainCertificate) {
	expect := expect(test)

	expect.NotNil("importSSLDomainCertificate", request)
}

/*
 * Test responses.
 */

var listSSLDomainCertificatesInNetworkDomainTestResponse = `
{
    "sslDomainCertificate": [{
            "id": "d9c7168d-7700-4063-b657-c41084246ab7",
            "datacenterId": "NA10",
            "networkDomainId": "14dbfacf-0eec-4964-a0c2-ff3f7390246b",
            "name": "My Domain Certificate",
            "description": "An SSL Cert for website x",
            "state": "NORMAL",
            "createTime": "2017-07-20 T16: 56: 08.000Z",
            "expiryTime": "2019-07-20 T16: 56: 08.000Z"
        },
        {
            "id": "5b8c3c29-3bb1-4af9-b85f-bb127d04a024",
            "datacenterId": "NA10",
            "networkDomainId": "14dbfacf-0eec-4964-a0c2-ff3f7390246b",
            "name": "My Other Domain Certificate",
            "description": "An SSL Cert for website y",
            "state": "NORMAL",
            "createTime": "2017-07-21T09:19:08.000Z",
            "expiryTime": "2019-07-31T23:59:59.000Z"
        }
    ],
    "pageNumber": 1,
    "pageCount": 2,
    "totalCount": 2,
    "pageSize": 250
}
 `

func verifyListSSLDomainCertificatesInNetworkDomainTestResponse(test *testing.T, response *SSLDomainCertificates) {
	expect := expect(test)

	expect.NotNil("SSLDomainCertificates", response)

	expect.EqualsInt("SSLDomainCertificates.PageNumber", 1, response.PageNumber)
	expect.EqualsInt("SSLDomainCertificates.PageCount", 2, response.PageCount)
	expect.EqualsInt("SSLDomainCertificates.TotalCount", 2, response.TotalCount)
	expect.EqualsInt("SSLDomainCertificates.PageSize", 250, response.PageSize)

	expect.EqualsInt("SSLDomainCertificates.Items.Length", 2, len(response.Items))

	certificate1 := response.Items[0]
	expect.EqualsString("SSLDomainCertificates[0].ID", "d9c7168d-7700-4063-b657-c41084246ab7", certificate1.ID)
	expect.EqualsString("SSLDomainCertificates[0].Name", "My Domain Certificate", certificate1.Name)
	expect.EqualsString("SSLDomainCertificates[0].Description", "An SSL Cert for website x", certificate1.Description)
	expect.EqualsString("SSLDomainCertificates[0].DatacenterID", "NA10", certificate1.DatacenterID)
	expect.EqualsString("SSLDomainCertificates[0].NetworkDomainID", "14dbfacf-0eec-4964-a0c2-ff3f7390246b", certificate1.NetworkDomainID)
	expect.EqualsString("SSLDomainCertificates[0].State", ResourceStatusNormal, certificate1.State)

	certificate2 := response.Items[1]
	expect.EqualsString("SSLDomainCertificates[1].ID", "5b8c3c29-3bb1-4af9-b85f-bb127d04a024", certificate2.ID)
	expect.EqualsString("SSLDomainCertificates[1].Name", "My Other Domain Certificate", certificate2.Name)
	expect.EqualsString("SSLDomainCertificates[1].Description", "An SSL Cert for website y", certificate2.Description)
	expect.EqualsString("SSLDomainCertificates[1].DatacenterID", "NA10", certificate2.DatacenterID)
	expect.EqualsString("SSLDomainCertificates[1].NetworkDomainID", "14dbfacf-0eec-4964-a0c2-ff3f7390246b", certificate2.NetworkDomainID)
	expect.EqualsString("SSLDomainCertificates[1].State", ResourceStatusNormal, certificate2.State)
}

var importSSLDomainCertificateTestRespose = `
{
    "operation": "IMPORT_SSL_DOMAIN_CERTIFICATE",
    "responseCode": "OK",
    "message": "SSL Domain Certificate 'Test-Ssl-Domain-Certificate' has been imported to Network Domain 9ace48a3-1f00-4635-aacf-e719b23cbf6c.",
    "info": [{
        "name": "sslDomainCertificateId",
        "value": "9e6b496d-5261-4542-91aa-b50c7f569c54"
    }],
    "requestId": "na9_20170821T074626030-0400_7e9fffe7-190b-46f2-9107-9d52fe57d0ad"
}
`

func verifyImportSSLDomainCertificateTestRespose(test *testing.T, request *importSSLDomainCertificate) {
	expect := expect(test)

	expect.NotNil("importSSLDomainCertificate", request)
	expect.EqualsString("importSSLDomainCertificate.Certificate", "-----BEGIN CERTIFICATE-----\nMIIFwzC...truncated for documentation brevity purposes...0xRcT7WQYUtUxu\n-----END CERTIFICATE-----", request.Certificate)
	expect.EqualsString("importSSLDomainCertificate.Key", "-----BEGIN PRIVATE KEY-----\nBAQEFAASCBKgwggGkl...truncated for documentation brevity purposes...wBPNAUcZVF5umh6GP\n-----END PRIVATE KEY-----", request.Key)
}
