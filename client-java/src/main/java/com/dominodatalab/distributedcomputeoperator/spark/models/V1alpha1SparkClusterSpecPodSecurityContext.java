/*
 * Kubernetes
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: v1.18.2
 * 
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */


package com.dominodatalab.distributedcomputeoperator.spark.models;

import java.util.Objects;
import java.util.Arrays;
import com.dominodatalab.distributedcomputeoperator.spark.models.V1alpha1SparkClusterSpecMasterLifecyclePostStartHttpGetHttpHeaders;
import com.dominodatalab.distributedcomputeoperator.spark.models.V1alpha1SparkClusterSpecMasterSecurityContextSeLinuxOptions;
import com.dominodatalab.distributedcomputeoperator.spark.models.V1alpha1SparkClusterSpecMasterSecurityContextSeccompProfile;
import com.dominodatalab.distributedcomputeoperator.spark.models.V1alpha1SparkClusterSpecMasterSecurityContextWindowsOptions;
import com.google.gson.TypeAdapter;
import com.google.gson.annotations.JsonAdapter;
import com.google.gson.annotations.SerializedName;
import com.google.gson.stream.JsonReader;
import com.google.gson.stream.JsonWriter;
import io.swagger.annotations.ApiModel;
import io.swagger.annotations.ApiModelProperty;
import java.io.IOException;
import java.util.ArrayList;
import java.util.List;

/**
 * V1alpha1SparkClusterSpecPodSecurityContext
 */
@javax.annotation.Generated(value = "org.openapitools.codegen.languages.JavaClientCodegen", date = "2021-10-04T16:38:59.380Z[Etc/UTC]")
public class V1alpha1SparkClusterSpecPodSecurityContext {
  public static final String SERIALIZED_NAME_FS_GROUP = "fsGroup";
  @SerializedName(SERIALIZED_NAME_FS_GROUP)
  private Long fsGroup;

  public static final String SERIALIZED_NAME_FS_GROUP_CHANGE_POLICY = "fsGroupChangePolicy";
  @SerializedName(SERIALIZED_NAME_FS_GROUP_CHANGE_POLICY)
  private String fsGroupChangePolicy;

  public static final String SERIALIZED_NAME_RUN_AS_GROUP = "runAsGroup";
  @SerializedName(SERIALIZED_NAME_RUN_AS_GROUP)
  private Long runAsGroup;

  public static final String SERIALIZED_NAME_RUN_AS_NON_ROOT = "runAsNonRoot";
  @SerializedName(SERIALIZED_NAME_RUN_AS_NON_ROOT)
  private Boolean runAsNonRoot;

  public static final String SERIALIZED_NAME_RUN_AS_USER = "runAsUser";
  @SerializedName(SERIALIZED_NAME_RUN_AS_USER)
  private Long runAsUser;

  public static final String SERIALIZED_NAME_SE_LINUX_OPTIONS = "seLinuxOptions";
  @SerializedName(SERIALIZED_NAME_SE_LINUX_OPTIONS)
  private V1alpha1SparkClusterSpecMasterSecurityContextSeLinuxOptions seLinuxOptions;

  public static final String SERIALIZED_NAME_SECCOMP_PROFILE = "seccompProfile";
  @SerializedName(SERIALIZED_NAME_SECCOMP_PROFILE)
  private V1alpha1SparkClusterSpecMasterSecurityContextSeccompProfile seccompProfile;

  public static final String SERIALIZED_NAME_SUPPLEMENTAL_GROUPS = "supplementalGroups";
  @SerializedName(SERIALIZED_NAME_SUPPLEMENTAL_GROUPS)
  private List<Long> supplementalGroups = null;

  public static final String SERIALIZED_NAME_SYSCTLS = "sysctls";
  @SerializedName(SERIALIZED_NAME_SYSCTLS)
  private List<V1alpha1SparkClusterSpecMasterLifecyclePostStartHttpGetHttpHeaders> sysctls = null;

  public static final String SERIALIZED_NAME_WINDOWS_OPTIONS = "windowsOptions";
  @SerializedName(SERIALIZED_NAME_WINDOWS_OPTIONS)
  private V1alpha1SparkClusterSpecMasterSecurityContextWindowsOptions windowsOptions;


  public V1alpha1SparkClusterSpecPodSecurityContext fsGroup(Long fsGroup) {
    
    this.fsGroup = fsGroup;
    return this;
  }

   /**
   * Get fsGroup
   * @return fsGroup
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public Long getFsGroup() {
    return fsGroup;
  }


  public void setFsGroup(Long fsGroup) {
    this.fsGroup = fsGroup;
  }


  public V1alpha1SparkClusterSpecPodSecurityContext fsGroupChangePolicy(String fsGroupChangePolicy) {
    
    this.fsGroupChangePolicy = fsGroupChangePolicy;
    return this;
  }

   /**
   * Get fsGroupChangePolicy
   * @return fsGroupChangePolicy
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public String getFsGroupChangePolicy() {
    return fsGroupChangePolicy;
  }


  public void setFsGroupChangePolicy(String fsGroupChangePolicy) {
    this.fsGroupChangePolicy = fsGroupChangePolicy;
  }


  public V1alpha1SparkClusterSpecPodSecurityContext runAsGroup(Long runAsGroup) {
    
    this.runAsGroup = runAsGroup;
    return this;
  }

   /**
   * Get runAsGroup
   * @return runAsGroup
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public Long getRunAsGroup() {
    return runAsGroup;
  }


  public void setRunAsGroup(Long runAsGroup) {
    this.runAsGroup = runAsGroup;
  }


  public V1alpha1SparkClusterSpecPodSecurityContext runAsNonRoot(Boolean runAsNonRoot) {
    
    this.runAsNonRoot = runAsNonRoot;
    return this;
  }

   /**
   * Get runAsNonRoot
   * @return runAsNonRoot
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public Boolean getRunAsNonRoot() {
    return runAsNonRoot;
  }


  public void setRunAsNonRoot(Boolean runAsNonRoot) {
    this.runAsNonRoot = runAsNonRoot;
  }


  public V1alpha1SparkClusterSpecPodSecurityContext runAsUser(Long runAsUser) {
    
    this.runAsUser = runAsUser;
    return this;
  }

   /**
   * Get runAsUser
   * @return runAsUser
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public Long getRunAsUser() {
    return runAsUser;
  }


  public void setRunAsUser(Long runAsUser) {
    this.runAsUser = runAsUser;
  }


  public V1alpha1SparkClusterSpecPodSecurityContext seLinuxOptions(V1alpha1SparkClusterSpecMasterSecurityContextSeLinuxOptions seLinuxOptions) {
    
    this.seLinuxOptions = seLinuxOptions;
    return this;
  }

   /**
   * Get seLinuxOptions
   * @return seLinuxOptions
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1SparkClusterSpecMasterSecurityContextSeLinuxOptions getSeLinuxOptions() {
    return seLinuxOptions;
  }


  public void setSeLinuxOptions(V1alpha1SparkClusterSpecMasterSecurityContextSeLinuxOptions seLinuxOptions) {
    this.seLinuxOptions = seLinuxOptions;
  }


  public V1alpha1SparkClusterSpecPodSecurityContext seccompProfile(V1alpha1SparkClusterSpecMasterSecurityContextSeccompProfile seccompProfile) {
    
    this.seccompProfile = seccompProfile;
    return this;
  }

   /**
   * Get seccompProfile
   * @return seccompProfile
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1SparkClusterSpecMasterSecurityContextSeccompProfile getSeccompProfile() {
    return seccompProfile;
  }


  public void setSeccompProfile(V1alpha1SparkClusterSpecMasterSecurityContextSeccompProfile seccompProfile) {
    this.seccompProfile = seccompProfile;
  }


  public V1alpha1SparkClusterSpecPodSecurityContext supplementalGroups(List<Long> supplementalGroups) {
    
    this.supplementalGroups = supplementalGroups;
    return this;
  }

  public V1alpha1SparkClusterSpecPodSecurityContext addSupplementalGroupsItem(Long supplementalGroupsItem) {
    if (this.supplementalGroups == null) {
      this.supplementalGroups = new ArrayList<>();
    }
    this.supplementalGroups.add(supplementalGroupsItem);
    return this;
  }

   /**
   * Get supplementalGroups
   * @return supplementalGroups
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public List<Long> getSupplementalGroups() {
    return supplementalGroups;
  }


  public void setSupplementalGroups(List<Long> supplementalGroups) {
    this.supplementalGroups = supplementalGroups;
  }


  public V1alpha1SparkClusterSpecPodSecurityContext sysctls(List<V1alpha1SparkClusterSpecMasterLifecyclePostStartHttpGetHttpHeaders> sysctls) {
    
    this.sysctls = sysctls;
    return this;
  }

  public V1alpha1SparkClusterSpecPodSecurityContext addSysctlsItem(V1alpha1SparkClusterSpecMasterLifecyclePostStartHttpGetHttpHeaders sysctlsItem) {
    if (this.sysctls == null) {
      this.sysctls = new ArrayList<>();
    }
    this.sysctls.add(sysctlsItem);
    return this;
  }

   /**
   * Get sysctls
   * @return sysctls
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public List<V1alpha1SparkClusterSpecMasterLifecyclePostStartHttpGetHttpHeaders> getSysctls() {
    return sysctls;
  }


  public void setSysctls(List<V1alpha1SparkClusterSpecMasterLifecyclePostStartHttpGetHttpHeaders> sysctls) {
    this.sysctls = sysctls;
  }


  public V1alpha1SparkClusterSpecPodSecurityContext windowsOptions(V1alpha1SparkClusterSpecMasterSecurityContextWindowsOptions windowsOptions) {
    
    this.windowsOptions = windowsOptions;
    return this;
  }

   /**
   * Get windowsOptions
   * @return windowsOptions
  **/
  @javax.annotation.Nullable
  @ApiModelProperty(value = "")

  public V1alpha1SparkClusterSpecMasterSecurityContextWindowsOptions getWindowsOptions() {
    return windowsOptions;
  }


  public void setWindowsOptions(V1alpha1SparkClusterSpecMasterSecurityContextWindowsOptions windowsOptions) {
    this.windowsOptions = windowsOptions;
  }


  @Override
  public boolean equals(Object o) {
    if (this == o) {
      return true;
    }
    if (o == null || getClass() != o.getClass()) {
      return false;
    }
    V1alpha1SparkClusterSpecPodSecurityContext v1alpha1SparkClusterSpecPodSecurityContext = (V1alpha1SparkClusterSpecPodSecurityContext) o;
    return Objects.equals(this.fsGroup, v1alpha1SparkClusterSpecPodSecurityContext.fsGroup) &&
        Objects.equals(this.fsGroupChangePolicy, v1alpha1SparkClusterSpecPodSecurityContext.fsGroupChangePolicy) &&
        Objects.equals(this.runAsGroup, v1alpha1SparkClusterSpecPodSecurityContext.runAsGroup) &&
        Objects.equals(this.runAsNonRoot, v1alpha1SparkClusterSpecPodSecurityContext.runAsNonRoot) &&
        Objects.equals(this.runAsUser, v1alpha1SparkClusterSpecPodSecurityContext.runAsUser) &&
        Objects.equals(this.seLinuxOptions, v1alpha1SparkClusterSpecPodSecurityContext.seLinuxOptions) &&
        Objects.equals(this.seccompProfile, v1alpha1SparkClusterSpecPodSecurityContext.seccompProfile) &&
        Objects.equals(this.supplementalGroups, v1alpha1SparkClusterSpecPodSecurityContext.supplementalGroups) &&
        Objects.equals(this.sysctls, v1alpha1SparkClusterSpecPodSecurityContext.sysctls) &&
        Objects.equals(this.windowsOptions, v1alpha1SparkClusterSpecPodSecurityContext.windowsOptions);
  }

  @Override
  public int hashCode() {
    return Objects.hash(fsGroup, fsGroupChangePolicy, runAsGroup, runAsNonRoot, runAsUser, seLinuxOptions, seccompProfile, supplementalGroups, sysctls, windowsOptions);
  }


  @Override
  public String toString() {
    StringBuilder sb = new StringBuilder();
    sb.append("class V1alpha1SparkClusterSpecPodSecurityContext {\n");
    sb.append("    fsGroup: ").append(toIndentedString(fsGroup)).append("\n");
    sb.append("    fsGroupChangePolicy: ").append(toIndentedString(fsGroupChangePolicy)).append("\n");
    sb.append("    runAsGroup: ").append(toIndentedString(runAsGroup)).append("\n");
    sb.append("    runAsNonRoot: ").append(toIndentedString(runAsNonRoot)).append("\n");
    sb.append("    runAsUser: ").append(toIndentedString(runAsUser)).append("\n");
    sb.append("    seLinuxOptions: ").append(toIndentedString(seLinuxOptions)).append("\n");
    sb.append("    seccompProfile: ").append(toIndentedString(seccompProfile)).append("\n");
    sb.append("    supplementalGroups: ").append(toIndentedString(supplementalGroups)).append("\n");
    sb.append("    sysctls: ").append(toIndentedString(sysctls)).append("\n");
    sb.append("    windowsOptions: ").append(toIndentedString(windowsOptions)).append("\n");
    sb.append("}");
    return sb.toString();
  }

  /**
   * Convert the given object to string with each line indented by 4 spaces
   * (except the first line).
   */
  private String toIndentedString(Object o) {
    if (o == null) {
      return "null";
    }
    return o.toString().replace("\n", "\n    ");
  }

}

