#import httplib, urllib
import os
import sys
import requests

def getRequiredEnvVar(name):
    envVar = os.environ.get(name)
    if envVar is None:
        print 'Required environment variable: ' + name + ' is missing. Exiting...'
        sys.exit(1)
    return envVar

def main():
    owner = getRequiredEnvVar('OWNER')
    package = getRequiredEnvVar('PACKAGE')
    productVersion = getRequiredEnvVar('PRODUCT_VERSION')
    githubToken = getRequiredEnvVar('GITHUB_TOKEN')

    # conn = httplib.HTTPSConnection("api.github.com")

    githubApiBaseUrl = 'https://api.github.com'
    #payloadCreateRelease = "{\"tag_name\": \"v1.0.15\"}"
    # urllib.urlencode({'tag_name': 'v1.0.15'})
    # conn.set_debuglevel(1)
    # conn.request("POST", "/repos/palette-software/insight-tester/releases", payloadCreateRelease, headers)
    # r1 = conn.getresponse()
    #
    # print r1.status, r1.reason

    headers = {'Authorization': 'token ' + githubToken}

    payload = {'tag_name': productVersion}


    r = requests.post(githubApiBaseUrl + '/repos/' + owner + '/' + package + '/releases', json=payload, headers=headers)
    print r.status_code, r.reason
    print r.text
    responseJson = r.json()
    if r.status_code == 200:
        releaseId = responseJson['id']
        if releaseId is None:
            # This is unexpected. There is supposed to be a Github release ID in the response JSON.
            sys.exit(1)

        # We are all good, let's print the ID of the newly created Github release
        print releaseId
        sys.exit(0)

    releaseAlreadyExists = False

    if r.status_code == 422:
        print responseJson['message']
        for error in responseJson['errors']:
            print error['code']
            if error['code'] == "already_exists":
                releaseAlreadyExists = True
                break
        if releaseAlreadyExists == False:
            sys.exit(2)

        r = requests.get(githubApiBaseUrl + '/repos/' + owner + '/' + package + '/releases', headers=headers)
        print 'Listing existing releases'
        print r.status_code, r.reason
        # print r.text

        for ver in r.json():
            if ver['tag_name'] == productVersion:
                releaseId = ver['id']
                if releaseId is None:
                    # This is unexpected. There is supposed to be a Github release ID in the response JSON.
                    sys.exit(1)

                # We are all good, let's print the ID of the existing Github release
                print releaseId
                sys.exit(0)

    # Github release ID was not found, and the response code was unexpected anyway
    print 'Response status code:' + r.status_code
    print 'Response message: ' + r.text
    sys.exit(3)

if __name__ == "__main__":
    main()
