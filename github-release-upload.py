import os
import sys
import requests

def getRequiredEnvVar(name):
    envVar = os.environ.get(name)
    if envVar is None:
        print 'Required environment variable: ' + name + ' is missing. Exiting...'
        sys.exit(1)
    return envVar

# Creates a new Github release and returns the ID of the newly created release. Although it returns the existing
# release ID, if the release already exits.
def main():
    owner = getRequiredEnvVar('OWNER')
    package = getRequiredEnvVar('PACKAGE')
    productVersion = getRequiredEnvVar('PRODUCT_VERSION')
    githubToken = getRequiredEnvVar('GITHUB_TOKEN')

    githubApiBaseUrl = 'https://api.github.com'
    headers = {'Authorization': 'token ' + githubToken}
    payload = {'tag_name': productVersion}

    r = requests.post(githubApiBaseUrl + '/repos/' + owner + '/' + package + '/releases', json=payload, headers=headers)
    # print r.status_code, r.reason
    # print r.text
    responseJson = r.json()
    if r.status_code == 200:
        releaseId = responseJson['id']
        if releaseId is None:
            # This is unexpected. There is supposed to be a Github release ID in the response JSON.
            print 'Error! Release ID was not found in JSON response!'
            sys.exit(1)

        # We are all good, let's print the ID of the newly created Github release
        print releaseId
        sys.exit(0)

    releaseAlreadyExists = False

    if r.status_code == 422:
        for error in responseJson['errors']:
            if error['code'] == "already_exists":
                releaseAlreadyExists = True
                break
        if releaseAlreadyExists == False:
            print 'Error! Release was expected to already exist, but it does not.'
            sys.exit(2)

        r = requests.get(githubApiBaseUrl + '/repos/' + owner + '/' + package + '/releases', headers=headers)

        for ver in r.json():
            if ver['tag_name'] == productVersion:
                releaseId = ver['id']
                if releaseId is None:
                    # This is unexpected. There is supposed to be a Github release ID in the response JSON.
                    print 'Error! Release ID was not found in the existing Github release!'
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
